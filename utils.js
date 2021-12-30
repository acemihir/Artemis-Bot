const config = require('./config');
const Filter = require('bad-words');
const fs = require('fs');

const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args));

module.exports.createId = function (prefix) {
    const chars = config.ids.chars;

    let result = '';
    for (let i = 0; i < config.ids.charLength; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }

    return prefix + result;
};

const wordFilter = new Filter();

module.exports.filterText = function (text) {
    return wordFilter.clean(text);
};

const addZeroBefore = function (n) {
    return (n < 10 ? '0' : '') + n;
};

module.exports.printLog = function (txt, type, shard = null) {
    const d = new Date();

    let func, pref;

    if (type === 'DEBUG') {
        func = console.debug;
        pref = shard === null ? '[DEBUG]' : `[DEBUG-${shard}]`;
    } else if (type === 'WARN') {
        func = console.warn;
        pref = shard === null ? '[WARN ]' : `[WARN-${shard} ]`;
    } else if (type === 'ERROR') {
        func = console.error;
        pref = shard === null ? '[ERROR]' : `[ERROR-${shard}]`;
    } else {
        func = console.info;
        pref = shard === null ? '[INFO ]' : `[INFO-${shard} ]`;
    }

    func(`${pref} ${addZeroBefore(d.getHours())}:${addZeroBefore(d.getSeconds())} > ${txt}`);
};

const getMember = async (gId, mId) => {
    const res = await fetch(`https://discord.com/api/v8/guilds/${gId}/members/${mId}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bot ' + config.botToken
        }
    }).catch(ex => this.printLog(ex, 'ERROR'));

    const body = await res.json();
    return body;  
};

const getRoles = async (gId) => {
    const res = await fetch(`https://discord.com/api/v8/guilds/${gId}/roles`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bot ' + config.botToken
        }
    }).catch(ex => this.printLog(ex, 'ERROR'));

    const body = await res.json();
    return body;    
};

module.exports.isAdmin = async function (guildId, memberId) {
    const roleIds = (await getMember(guildId, memberId)).roles;
    const roles = await getRoles(guildId);

    for (let i = 0; i < roles.length; i++) {
        if (roleIds.includes(roles[i].id)) {
            const perm = parseInt(roles[i].permissions);
            if ((perm & 0x8) == 8) {
                return true;
            }    
        }
    }

    return false;
};

module.exports.getFiles = (path) => {
    return new Promise(resolve => {
        fs.readdir(path, (ex, files) => {
            resolve(files.filter(f => f.endsWith('.js')));
        });
    });
};