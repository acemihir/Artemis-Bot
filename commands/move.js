// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
    .setName('move')
    .setDescription('Move a suggestion/report.')
    .setDefaultPermission(false)  
    .addStringOption(opt => opt.setName('id').setDescription('The ID of the suggestion/report.').setRequired(true))
    .addChannelOption(opt => opt.setName('channel').setDescription('The channel where the message should be moved to.').setRequired(true))

const execute = async function(client, interaction) {
    await interaction.reply('to be done')
}

// ================================
module.exports.command = {
    isPremium: false,
    privileged: true,

    data: data,
    execute: execute
}