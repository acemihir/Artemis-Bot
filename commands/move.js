// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
    .setName('move')
    .setDescription('Move a suggestion/report/poll.')
    .setDefaultPermission(false)  

    // Suggestions
    .addSubcommand(scmd =>
        scmd.setName('suggestions').setDescription('Move the suggestion message to another channel.')
            .addStringOption(opt => opt.setName('id').setDescription('The ID of the suggestion.'))
            .addChannelOption(opt => opt.setName('channel').setDescription('The channel where the message should be moved to.'))
    )

    // Reports
    .addSubcommand(scmd =>
        scmd.setName('reports').setDescription('Move the suggestion report to another channel.')
            .addStringOption(opt => opt.setName('id').setDescription('The ID of the report.'))
            .addChannelOption(opt => opt.setName('channel').setDescription('The channel where the message should be moved to.'))
    )

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