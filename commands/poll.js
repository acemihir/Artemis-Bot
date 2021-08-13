// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId } = require('../utils')

// ================================
const data = new SlashCommandBuilder()
    .setName('poll')
    .setDescription('Create a poll.')

const execute = async function(client, interaction) {
    await interaction.reply('to be done')
    const repId = createId('p_')
}

// ================================
module.exports.command = {
    isPremium: false,
    permLevel: 1,

    data: data,
    execute: execute
}