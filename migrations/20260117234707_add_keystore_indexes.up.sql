CREATE INDEX keystore_user_status_idx
ON keystore (user_id, status);

CREATE INDEX keystore_user_pkey_status_idx
ON keystore (user_id, p_key, status);

CREATE INDEX keystore_user_pkey_skey_status_idx
ON keystore (user_id, p_key, s_key, status);