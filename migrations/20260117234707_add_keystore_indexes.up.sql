CREATE INDEX keystore_client_status_idx
ON keystore (client_id, status);

CREATE INDEX keystore_client_pkey_status_idx
ON keystore (client_id, p_key, status);

CREATE INDEX keystore_client_pkey_skey_status_idx
ON keystore (client_id, p_key, s_key, status);