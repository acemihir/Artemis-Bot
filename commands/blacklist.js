// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
    .setName('blacklist')
    .setDescription('Blacklist a user from creating suggestions/reports.')
    .addUserOption(opt => opt.setName('member').setDescription('The member must be blacklisted').setRequired(true))
    .addStringOption(opt =>
        opt.setName('type')
            .setDescription('Where should the user be blacklisted for?')
            .addChoices([
                ['Suggestions', 'suggestions'],
                ['Reports', 'reports'],
                ['All', 'all']
            ]).setRequired(true))
    .setDefaultPermission(false)

const execute = async function(client, interaction) {
    interaction.reply('to be done')
}

// ================================
module.exports.command = {
    isPremium: true,
    privileged: true,

    data: data,
    execute: execute
}