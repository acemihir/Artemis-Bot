// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
    .setName('graph')
    .setDescription('Obtain a graph with information about suggestions/reports.')
    .addStringOption(opt =>
        opt.setName('type')
            .setDescription('The specific graph you want to see.')
            .addChoice('Suggestions', 'suggestions')
            .addChoice('Reports', 'reports'))

const execute = function(_client, interaction) {
    interaction.reply('to be done')
}

// ================================
module.exports.command = {
    isPremium: true,
    privileged: false,

    data: data,
    execute: execute
}