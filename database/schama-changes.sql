ALTER TABLE users ADD CONSTRAINT fk_user_coach FOREIGN KEY (coach_id) REFERENCES coach(id) ON DELETE SET NULL;
ALTER TABLE coach ADD CONSTRAINT fk_coach_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- after run all the tables
CREATE INDEX idx_workouts_user_id ON workouts(user_id);
CREATE INDEX idx_workout_exercises_workout_id ON workout_exercises(workout_id);
CREATE INDEX idx_workout_log_user_id ON workout_log(user_id);
CREATE INDEX idx_workout_log_performed_at ON workout_log(performed_at);

ALTER TABLE users ADD COLUMN token TEXT;
ALTER TABLE users ADD COLUMN refresh_token TEXT;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'user';
UPDATE users SET role = 'user' WHERE role = 'athlete';



