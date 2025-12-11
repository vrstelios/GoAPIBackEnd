CREATE TABLE users (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    password    TEXT NOT NULL,
    email       TEXT NOT NULL,
    role        TEXT NOT NULL DEFAULT 'athlete',
    coach_id    UUID NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coach_id) REFERENCES coach(id) ON DELETE SET NULL);

CREATE TABLE coach (
    id       UUID PRIMARY KEY,
    name     TEXT NOT NULL,
    user_id  UUID NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);

CREATE TABLE exercises (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    category    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE TABLE workouts (
    id           UUID PRIMARY KEY,
    user_id      UUID NOT NULL,
    name         TEXT NOT NULL,
    notes        TEXT,
    scheduled_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);

CREATE TABLE workout_exercises (
    id            UUID PRIMARY KEY,
    workout_id    UUID NOT NULL,
    exercise_id   UUID NOT NULL,
    sets          INT,
    reps          INT,
    weight        DECIMAL,
    rest_seconds  TEXT,
    exercise_order INT NOT NULL,
    UNIQUE(workout_id, exercise_id, exercise_order),
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE);

CREATE TABLE workout_log (
    id           UUID PRIMARY KEY,
    workout_id   UUID NOT NULL,
    user_id      UUID NOT NULL,
    performed_at TIMESTAMPTZ,
    duration     DECIMAL,
    total_volume TEXT,
    notes        TEXT,
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);

CREATE INDEX idx_workouts_user_id ON workouts(user_id);
CREATE INDEX idx_workout_exercises_workout_id ON workout_exercises(workout_id);
CREATE INDEX idx_workout_log_user_id ON workout_log(user_id);
CREATE INDEX idx_workout_log_performed_at ON workout_log(performed_at);


