alter table device_ir_remote_virtual_keys
    add column properties jsonb default '{}'::jsonb;