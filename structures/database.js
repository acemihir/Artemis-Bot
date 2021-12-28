const { Pool } = require('pg');
const config = require('../config');

const pool = new Pool({
    host: config.database.hostname,
    user: config.database.username,
    password: config.database.password,
    database: config.database.database,
    port: config.database.port,
    idleTimeoutMillis: 0,
    connectionTimeoutMillis: 0,
    max: 25
});

module.exports.runQuery = async function (query, params) {
    const client = await pool.connect();
    let result;
    try {
        result = !params ? await client.query(query) : await client.query(query, params);
    } catch (ex) {
        console.error(ex);
    } finally {
        client.release();
    }
    return result;
};