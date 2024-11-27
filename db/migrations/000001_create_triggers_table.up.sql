CREATE TABLE IF NOT EXISTS triggers (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    schedule_hour INTEGER,
    schedule_minute INTEGER,
    schedule_action TEXT,
    flic_button_mac TEXT,
    flic_button_on_single_click TEXT,
    flic_button_on_double_click TEXT,
    flic_button_on_hold TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);