"""Run: python3 -m unittest discover -s <recorder>/scripts -p 'test_*.py'."""
import copy
import json
import sqlite3
import subprocess
import sys
import tempfile
import unittest
from datetime import datetime, timedelta
from pathlib import Path
from unittest.mock import patch

import quiz


class QuizTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = Path(self.tmp.name)
        (self.root / 'model.py').write_text('required = True\n')
        self.path = self.root / 'quiz.sqlite3'
        self.db = quiz.connect(self.path)
        self.store = quiz.Store(self.db, self.root, 'repo-1')
        self.raw = {
            'scope': ['model.py'], 'prompt': 'What is guaranteed?',
            'choices': [{'id': str(i), 'text': f'Choice {i}', 'explanation': f'Reason {i}'} for i in range(4)],
            'correct_id': '2', 'explanation': 'Implemented guarantee', 'review_point': 'Check the layer',
            'evidence': [{'path': 'model.py', 'locator': 'required'}],
            'documentation_suggestion': 'Align the documentation.'}
        self.qid = self.add(self.raw)

    def tearDown(self):
        self.db.close()
        self.tmp.cleanup()

    def add(self, raw):
        with self.store.transaction():
            return self.store.add(quiz.question_input(self.root, raw))['id']

    def answer(self, presentation, correct=True):
        order = presentation['presentation']['order']
        choice = 'ABCD'[order.index('2')] if correct else 'ABCD'[order.index('0')]
        return self.store.answer(presentation, choice)

    def make_due(self):
        self.db.execute("UPDATE progress SET due_at='2000-01-01T00:00:00.000000+00:00' WHERE question_id=?", (self.qid,))

    def test_intervals_mastery_and_duplicate_answer(self):
        for streak, days in ((1, 1), (2, 7), (3, None)):
            self.make_due()
            presentation = self.store.present(self.qid, False)
            before = datetime.fromisoformat(quiz.now())
            result = self.answer(presentation)
            self.assertEqual(result['streak'], streak)
            if days:
                delta = datetime.fromisoformat(result['due_at']) - before
                self.assertAlmostEqual(delta.total_seconds(), timedelta(days=days).total_seconds(), delta=1)
                self.assertEqual(self.store.select(['.'], 1, False, [])['questions'], [])
            else:
                self.assertTrue(result['mastered'])
                self.assertIsNone(result['due_at'])
            with self.assertRaisesRegex(ValueError, 'already answered'):
                self.answer(presentation)
        self.assertEqual(self.store.get(self.qid)['presentations'], 3)
        with self.assertRaises(ValueError):
            self.store.present(self.qid, False)

    def test_wrong_answer_resets_streak(self):
        self.answer(self.store.present(self.qid, False))
        self.make_due()
        before = datetime.fromisoformat(quiz.now())
        result = self.answer(self.store.present(self.qid, False), False)
        self.assertEqual(result['streak'], 0)
        self.assertAlmostEqual((datetime.fromisoformat(result['due_at']) - before).total_seconds(), 86400, delta=1)

    def test_practice_does_not_write_or_consume_first_presentation(self):
        before = list(self.db.iterdump())
        presentation = self.store.present(self.qid, True)
        self.assertTrue(self.answer(presentation)['correct'])
        self.assertFalse(self.answer(presentation)['saved'])
        self.assertEqual(list(self.db.iterdump()), before)
        self.assertEqual(len(self.store.select(['.'], 1, True, [])['questions']), 1)
        self.assertEqual(self.store.present(self.qid, False)['display'], '初出')

    def test_unanswered_presentation_only_increases_count(self):
        self.store.present(self.qid, False)
        row = self.store.get(self.qid)
        self.assertEqual((row['presentations'], row['streak'], row['due_at']), (1, 0, None))
        self.assertEqual(self.store.select(['.'], 1, True, [])['questions'], [])
        self.assertEqual(self.store.present(self.qid, False)['display'], '今回で2回目')

    def test_answer_uses_saved_order_and_snapshot(self):
        presentation = self.store.present(self.qid, False)
        original = copy.deepcopy(presentation)
        presentation['presentation']['order'].reverse()
        saved = original['presentation']['order']
        result = self.store.answer(presentation, 'ABCD'[saved.index('2')])
        self.assertTrue(result['correct'])
        attempt = self.db.execute('SELECT * FROM attempts').fetchone()
        self.assertEqual(json.loads(attempt['snapshot'])['prompt'], self.raw['prompt'])

    def test_stale_evidence_blocks_selection_presentation_and_answer(self):
        presentation = self.store.present(self.qid, False)
        (self.root / 'model.py').write_text('required = False\n')
        result = self.store.select(['.'], 2, False, [])
        self.assertEqual(result['questions'], [])
        self.assertEqual(result['needs_review'], [self.qid])
        with self.assertRaisesRegex(ValueError, 'evidence changed'):
            self.store.present(self.qid, False)
        with self.assertRaises(ValueError):
            self.answer(presentation)
        self.assertEqual(self.store.get(self.qid)['streak'], 0)

    def command(self, name, **kwargs):
        return quiz.run(self.store, type('Args', (), {'command': name, **kwargs})())

    def test_revision_preserves_old_history_and_resets_new(self):
        self.answer(self.store.present(self.qid, False))
        revised = copy.deepcopy(self.raw)
        revised['prompt'] = 'Updated guarantee?'
        path = self.root / 'input.json'
        path.write_text(json.dumps(revised))
        result = self.command('revise', id=self.qid, input=str(path))
        self.assertEqual(self.store.get(self.qid)['status'], 'superseded')
        self.assertEqual(self.store.get(result['id'])['presentations'], 0)
        self.assertEqual(len(self.command('history', id=self.qid, limit=20)), 1)
        with self.assertRaises(ValueError):
            self.command('reset', id=self.qid)

    def test_duplicate_revision_leaves_original_active(self):
        path = self.root / 'input.json'
        path.write_text(json.dumps(self.raw))
        with self.assertRaises(ValueError):
            self.command('revise', id=self.qid, input=str(path))
        self.assertEqual(self.store.get(self.qid)['status'], 'active')

    def test_delete_reset_and_old_presentation(self):
        presentation = self.store.present(self.qid, False)
        self.command('delete', id=self.qid)
        self.assertEqual(self.store.select(['.'], 1, False, [])['questions'], [])
        self.command('reset', id=self.qid)
        self.assertEqual(self.store.get(self.qid)['presentations'], 1)
        with self.assertRaisesRegex(ValueError, 'reset|changed'):
            self.answer(presentation)

    def test_prior_presentation_cannot_change_newer_progress(self):
        older = self.store.present(self.qid, False)
        latest = self.store.present(self.qid, False)
        self.answer(latest)
        with self.assertRaisesRegex(ValueError, 'newer'):
            self.answer(older)
        self.assertEqual(self.store.get(self.qid)['streak'], 1)

    def test_selection_scope_order_exclusion_and_repository_isolation(self):
        self.store.present(self.qid, False)  # Prior, unanswered problem remains due.
        raw = copy.deepcopy(self.raw)
        raw['prompt'] = 'New problem'
        other = self.add(raw)
        ids = [r['id'] for r in self.store.select(['model.py'], 3, False, [])['questions']]
        self.assertEqual(ids, [self.qid, other])
        self.assertEqual(self.store.select(['unrelated'], 3, False, [])['shortfall'], 3)
        self.assertEqual([r['id'] for r in self.store.select(['.'], 3, True, [])['questions']], [other])
        self.assertEqual(self.store.select(['.'], 3, False, [self.qid, other])['shortfall'], 3)
        isolated = quiz.Store(self.db, self.root, 'repo-2')
        self.assertEqual(isolated.rows(), [])
        with self.assertRaises(ValueError):
            isolated.get(self.qid)

    def test_input_constraints_and_evidence_escape(self):
        for key, value in [('choices', self.raw['choices'][:3]), ('correct_id', 'missing'),
                           ('evidence', []), ('scope', ['../outside'])]:
            invalid = {**self.raw, key: value}
            with self.assertRaises(ValueError):
                quiz.question_input(self.root, invalid)
        (self.root / 'escape').symlink_to('/etc/hosts')
        with self.assertRaises(ValueError):
            quiz.file_hash(self.root, 'escape')

    def test_future_and_unversioned_databases_are_refused(self):
        self.db.execute('PRAGMA user_version=99')
        with self.assertRaisesRegex(ValueError, 'unsupported'):
            quiz.connect(self.path)
        self.assertEqual(self.db.execute('PRAGMA user_version').fetchone()[0], 99)
        self.db.execute('PRAGMA user_version=0')
        with self.assertRaisesRegex(ValueError, 'unversioned'):
            quiz.connect(self.path)

    def test_migration_backup_success_and_rollback(self):
        migration_dir = self.root / 'migrations'
        migration_dir.mkdir()
        sql = migration_dir / '002.sql'
        sql.write_text('CREATE TABLE extra(id INTEGER); PRAGMA user_version=2;')
        with patch.object(quiz, 'VERSION', 2), patch.object(quiz, '__file__', str(self.root / 'quiz.py')):
            quiz.migrate(self.db, self.path)
        self.assertEqual(self.db.execute('PRAGMA user_version').fetchone()[0], 2)
        backups = list(self.root.glob('*.backup-*'))
        self.assertEqual(len(backups), 1)
        with sqlite3.connect(backups[0]) as backup:
            self.assertEqual(backup.execute('SELECT count(*) FROM questions').fetchone()[0], 1)
            self.assertEqual(backup.execute('PRAGMA user_version').fetchone()[0], 1)
        (migration_dir / '003.sql').write_text('CREATE TABLE partial(id); INVALID SQL; PRAGMA user_version=3;')
        with patch.object(quiz, 'VERSION', 3), patch.object(quiz, '__file__', str(self.root / 'quiz.py')):
            with self.assertRaises(sqlite3.Error):
                quiz.migrate(self.db, self.path)
        self.assertEqual(self.db.execute('PRAGMA user_version').fetchone()[0], 2)
        self.assertIsNone(self.db.execute("SELECT name FROM sqlite_master WHERE name='partial'").fetchone())

    def test_cli_end_to_end(self):
        script = Path(quiz.__file__).resolve()
        repo = script.parents[5]
        db_path = self.root / 'cli.sqlite3'
        raw = copy.deepcopy(self.raw)
        raw['scope'] = ['AGENTS.md']
        raw['evidence'] = [{'path': 'AGENTS.md', 'locator': 'project instructions'}]
        input_path = self.root / 'cli-input.json'
        input_path.write_text(json.dumps(raw))

        def cli(*args):
            result = subprocess.run([sys.executable, str(script), '--repo', str(repo), '--db', str(db_path), *args],
                                    text=True, capture_output=True, check=True)
            return json.loads(result.stdout)

        qid = cli('add', '--input', str(input_path))['id']
        self.assertEqual(cli('select', '--scope', '.', '--count', '1')['questions'][0]['id'], qid)
        presentation = cli('present', '--id', str(qid))
        self.assertNotIn('correct_id', presentation)
        presentation_path = self.root / 'presentation.json'
        presentation_path.write_text(json.dumps(presentation))
        correct = 'ABCD'[presentation['presentation']['order'].index('2')]
        self.assertTrue(cli('answer', '--presentation', str(presentation_path), '--choice', correct)['correct'])
        self.assertEqual(len(cli('history')), 1)
        self.assertEqual(cli('list')[0]['streak'], 1)

    def test_missing_migration_version_rolls_back(self):
        migration_dir = self.root / 'migrations'
        migration_dir.mkdir()
        (migration_dir / '002.sql').write_text('CREATE TABLE incomplete(id);')
        with patch.object(quiz, 'VERSION', 2), patch.object(quiz, '__file__', str(self.root / 'quiz.py')):
            with self.assertRaisesRegex(ValueError, 'expected schema version'):
                quiz.migrate(self.db, self.path)
        self.assertEqual(self.db.execute('PRAGMA user_version').fetchone()[0], 1)
        self.assertIsNone(self.db.execute("SELECT name FROM sqlite_master WHERE name='incomplete'").fetchone())


if __name__ == '__main__':
    unittest.main()
