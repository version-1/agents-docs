#!/usr/bin/env python3
"""Local codebase quizzes. All commands emit JSON; no third-party dependencies."""
import argparse
import hashlib
import json
import random
import sqlite3
import subprocess
import sys
import uuid
from contextlib import contextmanager
from datetime import datetime, timedelta, timezone
from pathlib import Path

VERSION = 2
DEFAULT_DB = Path.home() / '.local/share/quiz-code-base/quiz.sqlite3'


def encode(value):
    return json.dumps(value, ensure_ascii=False, sort_keys=True)


def now():
    return datetime.now(timezone.utc).isoformat(timespec='microseconds')


def digest(value):
    return hashlib.sha256(encode(value).encode()).hexdigest()


def git(repo, *args):
    return subprocess.check_output(['git', '-C', str(repo), *args], text=True,
                                   stderr=subprocess.PIPE).strip()


def repository(repo):
    root = Path(git(repo, 'rev-parse', '--show-toplevel')).resolve()
    common = Path(git(root, 'rev-parse', '--git-common-dir'))
    return root, str((root / common).resolve() if not common.is_absolute() else common.resolve())


def relative(value):
    path = Path(value)
    if path.is_absolute() or '..' in path.parts:
        raise ValueError('scope/evidence paths must be repository-relative without ..')
    return path.as_posix()


def file_hash(root, name):
    root = root.resolve()
    path = (root / relative(name)).resolve()
    if not path.is_relative_to(root):
        raise ValueError('evidence escapes repository')
    return hashlib.sha256(path.read_bytes()).hexdigest()


def nonempty(value, label):
    if not isinstance(value, str) or not value.strip():
        raise ValueError(label + ' must be a nonempty string')


def question_input(root, data):
    if not isinstance(data, dict):
        raise ValueError('question input must be an object')
    for key in ('prompt', 'explanation', 'review_point'):
        nonempty(data.get(key), key)
    choices = data.get('choices')
    if not isinstance(choices, list) or len(choices) != 4:
        raise ValueError('exactly four choices required')
    for choice in choices:
        if not isinstance(choice, dict):
            raise ValueError('each choice must be an object')
        for key in ('id', 'text', 'explanation'):
            nonempty(choice.get(key), 'choice.' + key)
    if len({c['id'] for c in choices}) != 4 or len({c['text'].strip() for c in choices}) != 4:
        raise ValueError('choice IDs and texts must be distinct')
    if data.get('correct_id') not in {c['id'] for c in choices}:
        raise ValueError('correct_id must identify exactly one choice')
    scopes = data.get('scope')
    if not isinstance(scopes, list) or not scopes:
        raise ValueError('scope must be a nonempty list of repository paths')
    scopes = sorted({relative(s) for s in scopes})
    evidence = data.get('evidence')
    if not isinstance(evidence, list) or not evidence:
        raise ValueError('at least one evidence file required')
    stamped = []
    for item in evidence:
        if not isinstance(item, dict):
            raise ValueError('each evidence entry must be an object')
        name = relative(item['path'])
        nonempty(item.get('locator'), 'evidence.locator')
        stamped.append({'path': name, 'locator': item['locator'], 'sha256': file_hash(root, name)})
    suggestion = data.get('documentation_suggestion', '')
    if not isinstance(suggestion, str):
        raise ValueError('documentation_suggestion must be a string')
    return {**{k: data[k] for k in ('prompt', 'explanation', 'review_point', 'correct_id')},
            'scope': scopes, 'choices': choices, 'evidence': stamped,
            'documentation_suggestion': suggestion}


def fingerprint(data):
    # Evidence stamps/locators can change without changing the question's meaning.
    return digest({k: v for k, v in data.items() if k != 'evidence'})


def stale(root, data):
    for item in data['evidence']:
        try:
            if file_hash(root, item['path']) != item['sha256']:
                return True
        except (OSError, ValueError):
            return True
    return False


def connect(path):
    path = Path(path).expanduser().resolve()
    path.parent.mkdir(parents=True, exist_ok=True)
    db = sqlite3.connect(path, timeout=30, isolation_level=None)
    db.row_factory = sqlite3.Row
    db.execute('PRAGMA foreign_keys=ON')
    try:
        migrate(db, path)
    except Exception:
        db.close()
        raise
    return db


def migrate(db, path):
    version = db.execute('PRAGMA user_version').fetchone()[0]
    if version > VERSION or version < 0:
        raise ValueError('unsupported database version: ' + str(version))
    if version == VERSION:
        return
    tables = db.execute("SELECT name FROM sqlite_master WHERE type='table'").fetchall()
    if version == 0 and tables:
        raise ValueError('unversioned nonempty database; automatic migration refused')
    # sqlite backup also captures committed WAL contents.
    if tables:
        backup = sqlite3.connect(str(path) + '.backup-' + uuid.uuid4().hex)
        try:
            db.backup(backup)
        finally:
            backup.close()
    migrations = Path(__file__).parent / 'migrations'
    scripts = [(migrations / f'{n:03}.sql').read_text() for n in range(version + 1, VERSION + 1)]
    try:
        db.executescript('BEGIN EXCLUSIVE;\n' + '\n'.join(scripts))
        if db.execute('PRAGMA user_version').fetchone()[0] != VERSION:
            raise ValueError('migration did not set the expected schema version')
        db.commit()
    except Exception:
        if db.in_transaction:
            db.rollback()
        raise


class Store:
    def __init__(self, db, root, identity):
        self.db, self.root, self.identity = db, root, identity

    @contextmanager
    def transaction(self):
        self.db.execute('BEGIN IMMEDIATE')
        try:
            yield
            self.db.commit()
        except Exception:
            self.db.rollback()
            raise

    def repo_id(self):
        self.db.execute('INSERT OR IGNORE INTO repositories(identity) VALUES (?)', (self.identity,))
        return self.db.execute('SELECT id FROM repositories WHERE identity=?', (self.identity,)).fetchone()[0]

    def rows(self):
        return self.db.execute('''SELECT q.*, p.presentations, p.streak, p.due_at, p.mastered, p.generation
            FROM questions q JOIN progress p ON p.question_id=q.id
            JOIN repositories r ON r.id=q.repository_id WHERE r.identity=? ORDER BY q.id''',
                               (self.identity,)).fetchall()

    def get(self, qid):
        for row in self.rows():
            if row['id'] == qid:
                return row
        raise ValueError('question not found in this repository')

    def add(self, data, replaces=None):
        fp = fingerprint(data)
        rid = self.repo_id()
        existing = self.db.execute("SELECT id FROM questions WHERE repository_id=? AND fingerprint=? AND status!='superseded' ORDER BY id LIMIT 1",
                                   (rid, fp)).fetchone()
        if existing:
            if replaces is not None:
                raise ValueError('revision duplicates an existing question; original preserved')
            return {'id': existing[0], 'created': False}
        cur = self.db.execute('''INSERT INTO questions(repository_id,content,fingerprint,replaces,created_at)
            VALUES (?,?,?,?,?)''', (rid, encode(data), fp, replaces, now()))
        qid = cur.lastrowid
        self.db.execute('INSERT INTO progress(question_id) VALUES (?)', (qid,))
        return {'id': qid, 'created': True}

    def summary(self, row):
        data = json.loads(row['content'])
        return {**{k: row[k] for k in ('id', 'status', 'presentations', 'streak', 'due_at', 'mastered')},
                'scope': data['scope'], 'prompt': data['prompt']}

    def eligible(self, row, timestamp):
        return (row['status'] == 'active' and not row['mastered']
                and (row['due_at'] is None or row['due_at'] <= timestamp))

    def select(self, scopes, count, new_only, exclude):
        if count < 1:
            raise ValueError('count must be positive')
        scopes = [relative(p) for p in scopes]

        def matches(data):
            return any(s == '.' or p == s or p.startswith(s + '/')
                       for p in data['scope'] for s in scopes)

        candidates = [r for r in self.rows() if self.eligible(r, now()) and r['id'] not in exclude
                      and (not new_only or r['presentations'] == 0) and matches(json.loads(r['content']))]
        candidates.sort(key=lambda r: (r['presentations'] == 0, r['due_at'] or '', r['id']))
        selected, review = [], []
        for row in candidates:
            if stale(self.root, json.loads(row['content'])):
                if len(review) < count:
                    review.append(row['id'])
            elif len(selected) < count:
                selected.append(self.summary(row))
        return {'questions': selected, 'needs_review': review, 'shortfall': count - len(selected)}

    def present(self, qid, practice):
        with self.transaction():
            row = self.get(qid)
            data = json.loads(row['content'])
            if not self.eligible(row, now()):
                raise ValueError('question is not due/active or is mastered')
            if stale(self.root, data):
                raise ValueError('evidence changed; analyzer review required')
            choices = [c['id'] for c in data['choices']]
            random.SystemRandom().shuffle(choices)
            ordinal = row['presentations'] + (not practice)
            ticket = {'question_id': qid, 'practice': practice, 'order': choices,
                      'revision': digest(data), 'repository': self.identity, 'ordinal': ordinal,
                      'generation': row['generation']}
            if not practice:
                ticket['attempt_id'] = uuid.uuid4().hex
                self.db.execute('UPDATE progress SET presentations=presentations+1 WHERE question_id=?', (qid,))
                self.db.execute('''INSERT INTO attempts(id,question_id,presented_at,ordinal,generation,snapshot,choices)
                    VALUES (?,?,?,?,?,?,?)''', (ticket['attempt_id'], qid, now(), ordinal,
                                              row['generation'], encode(data), encode(choices)))
            by_id = {c['id']: c['text'] for c in data['choices']}
            return {'presentation': ticket, 'prompt': data['prompt'],
                    'choices': {label: by_id[c] for label, c in zip('ABCD', choices)},
                    'display': ('練習・記録上の出題回数' + str(ordinal) + '回（更新なし）' if practice else
                                '初出' if ordinal == 1 else f'今回で{ordinal}回目')}

    def answer(self, ticket, choice):
        ticket = ticket.get('presentation', ticket)
        if ticket['repository'] != self.identity:
            raise ValueError('presentation belongs to a different repository')
        with self.transaction():
            row = self.get(ticket['question_id'])
            data = json.loads(row['content'])
            practice = ticket['practice']
            if not isinstance(practice, bool):
                raise ValueError('practice must be boolean')
            order = ticket['order']
            attempt = None
            if not practice:
                attempt = self.db.execute('SELECT * FROM attempts WHERE id=? AND question_id=?',
                                          (ticket['attempt_id'], row['id'])).fetchone()
                if not attempt:
                    raise ValueError('unknown presentation')
                order = json.loads(attempt['choices'])
                if attempt['answered_at'] is not None:
                    raise ValueError('presentation already answered; progress unchanged')
                if json.loads(attempt['snapshot']) != data:
                    raise ValueError('question changed since presentation')
                if attempt['generation'] != row['generation']:
                    raise ValueError('learning state reset since presentation')
            if ticket['generation'] != row['generation']:
                raise ValueError('learning state changed since presentation')
            if row['status'] != 'active' or stale(self.root, data) or digest(data) != ticket['revision']:
                raise ValueError('question/evidence changed; answer not recorded')
            if len(order) != 4 or set(order) != {c['id'] for c in data['choices']}:
                raise ValueError('invalid presentation order')
            # Only the most recently presented attempt may change progress. This prevents stale
            # parallel conversations from advancing a newer learning state.
            if not practice and attempt['ordinal'] != row['presentations']:
                raise ValueError('a newer presentation exists; answer not recorded')
            selected = order['ABCD'.index(choice)]
            correct = selected == data['correct_id']
            streak = min(row['streak'] + 1, 3) if correct else 0
            timestamp = now()
            due = None if streak == 3 else (datetime.fromisoformat(timestamp) + timedelta(
                days=7 if streak == 2 else 1)).isoformat(timespec='microseconds')
            by_id = {c['id']: c for c in data['choices']}
            feedback = {'correct': correct, 'correct_choice': 'ABCD'[order.index(data['correct_id'])],
                        'explanation': data['explanation'], 'selected_explanation': by_id[selected]['explanation'],
                        'evidence': data['evidence'], 'review_point': data['review_point'],
                        'documentation_suggestion': data['documentation_suggestion'], 'saved': not practice,
                        'streak': row['streak'] if practice else streak,
                        'due_at': row['due_at'] if practice else due,
                        'mastered': bool(row['mastered']) if practice else streak == 3}
            if not practice:
                self.db.execute('UPDATE progress SET streak=?,due_at=?,mastered=? WHERE question_id=?',
                                (streak, due, streak == 3, row['id']))
                self.db.execute('UPDATE attempts SET answered_at=?,selected=?,correct=?,feedback=? WHERE id=?',
                                (timestamp, selected, correct, encode(feedback), ticket['attempt_id']))
            return feedback


def parser():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--db', default=str(DEFAULT_DB))
    p.add_argument('--repo', default='.')
    sub = p.add_subparsers(dest='command', required=True)
    sub.add_parser('init')
    listing = sub.add_parser('list')
    listing.add_argument('--all', action='store_true')
    for name in ('inspect', 'delete', 'reset', 'refresh'):
        sub.add_parser(name).add_argument('--id', type=int, required=True)
    for name in ('add', 'revise'):
        child = sub.add_parser(name)
        child.add_argument('--input', required=True)
        if name == 'revise':
            child.add_argument('--id', type=int, required=True)
    select = sub.add_parser('select')
    select.add_argument('--scope', action='append', required=True)
    select.add_argument('--count', type=int, required=True)
    select.add_argument('--new-only', action='store_true')
    select.add_argument('--exclude', nargs='*', type=int, default=[])
    present = sub.add_parser('present')
    present.add_argument('--id', type=int, required=True)
    present.add_argument('--practice', action='store_true')
    answer = sub.add_parser('answer')
    answer.add_argument('--presentation', required=True)
    answer.add_argument('--choice', choices=list('ABCD'), required=True)
    history = sub.add_parser('history')
    history.add_argument('--id', type=int)
    history.add_argument('--limit', type=int, default=20)
    return p


def run(store, args):
    cmd = args.command
    if cmd == 'init':
        return {'schema_version': VERSION}
    if cmd == 'list':
        return [store.summary(r) for r in store.rows() if args.all or r['status'] == 'active']
    if cmd == 'inspect':
        row = store.get(args.id)
        return {**store.summary(row), 'content': json.loads(row['content'])}
    if cmd == 'select':
        return store.select(args.scope, args.count, args.new_only, args.exclude)
    if cmd == 'present':
        return store.present(args.id, args.practice)
    if cmd == 'answer':
        return store.answer(json.loads(Path(args.presentation).read_text()), args.choice)
    if cmd == 'history':
        if args.limit < 1:
            raise ValueError('limit must be positive')
        if args.id is not None:
            store.get(args.id)
        rows = store.db.execute('''SELECT a.* FROM attempts a JOIN questions q ON q.id=a.question_id
            JOIN repositories r ON r.id=q.repository_id WHERE r.identity=? AND (? IS NULL OR q.id=?)
            ORDER BY a.presented_at DESC LIMIT ?''', (store.identity, args.id, args.id, args.limit)).fetchall()
        return [{**dict(r), 'snapshot': json.loads(r['snapshot']), 'choices': json.loads(r['choices']),
                 'feedback': json.loads(r['feedback']) if r['feedback'] else None} for r in rows]
    data = question_input(store.root, json.loads(Path(args.input).read_text())) if cmd in ('add', 'revise') else None
    with store.transaction():
        if cmd == 'add':
            return store.add(data)
        row = store.get(args.id)
        if row['status'] == 'superseded':
            raise ValueError('superseded revision is immutable')
        if cmd == 'revise':
            if data == json.loads(row['content']):
                raise ValueError('revision is unchanged; use refresh for evidence revalidation')
            store.db.execute("UPDATE questions SET status='superseded' WHERE id=?", (args.id,))
            result = store.add(data, replaces=args.id)
            return result
        if cmd == 'delete':
            store.db.execute("UPDATE questions SET status='deleted' WHERE id=?", (args.id,))
        elif cmd == 'reset':
            store.db.execute("UPDATE questions SET status='active' WHERE id=?", (args.id,))
            store.db.execute('UPDATE progress SET streak=0,mastered=0,due_at=NULL,generation=generation+1 WHERE question_id=?', (args.id,))
        elif cmd == 'refresh':
            data = question_input(store.root, json.loads(row['content']))
            store.db.execute('UPDATE questions SET content=? WHERE id=?', (encode(data), args.id))
        return {'id': args.id, 'operation': cmd}


def main():
    args = parser().parse_args()
    db = None
    try:
        root, identity = repository(args.repo)
        db = connect(args.db)
        result = run(Store(db, root, identity), args)
        print(encode(result))
        return 0
    except (ValueError, KeyError, TypeError, OSError, sqlite3.Error, subprocess.CalledProcessError) as error:
        print(encode({'error': str(error)}), file=sys.stderr)
        return 1
    finally:
        if db is not None:
            db.close()


if __name__ == '__main__':
    sys.exit(main())
