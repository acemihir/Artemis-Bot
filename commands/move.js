// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
    .setName('move')
    .setDescription('Move a suggestion/report/poll.')

    // Suggestions
    .addSubcommand(scmd =>
        scmd.setName('suggestions').setDescription('Move the suggestion message to another channel.')
            .addStringOption(opt => opt.setName('Suggestion ID').setDescription('The ID of the suggestion.'))
            .addChannelOption(opt => opt.setName('Channel').setDescription('The channel where the message should be moved to.')
                .setRequired(true))
    )

    // Reports
    .addSubcommand(scmd =>
        scmd.setName('reports').setDescription('Move the suggestion report to another channel.')
            .addStringOption(opt => opt.setName('Report ID').setDescription('The ID of the report.'))
            .addChannelOption(opt => opt.setName('Channel').setDescription('The channel where the message should be moved to.')
                .setRequired(true))
    )

    // Polls
    .addSubcommand(scmd =>
        scmd.setName('polls').setDescription('Move the suggestion poll to another channel.')
            .addStringOption(opt => opt.setName('Poll ID').setDescription('The ID of the poll.'))
            .addChannelOption(opt => opt.setName('Channel').setDescription('The channel where the message should be moved to.')
                .setRequired(true))
    )

const execute = async function(client, interaction) {
    await interaction.reply('to be done')
}

// ================================
module.exports.command = {
    isPremium: false,
    permLevel: 1,

    data: data,
    execute: execute
}