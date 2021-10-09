const { botCache, getFromRedis } = require('../structures/cache')

module.exports = async function(client, interaction) {
    if (interaction.isCommand()) {
        // Check if the used command is actually stored in the botCache object
        if (botCache.commands.has(interaction.commandName)) {
            // Retrieve the command data from the botCache object
            const obj = botCache.commands.get(interaction.commandName)

            // Fetch the guild data from the cache
            const cachedData = await getFromRedis(interaction.guildId)

            // Check if the command is a premium command
            if (obj.isPremium) {
                // Check if the guild does not have premium
                if (!cachedData.premium) {
                    // Construct the row
                    const row = new MessageActionRow().addComponents(new MessageButton()
                        .setURL('https://github.com/jerskisnow/Suggestions/wiki/Donating')
                        .setLabel('Donating')
                        .setEmoji('💰')
                        .setStyle('LINK'))

                    // Construct the embed
                    const embed = new MessageEmbed()
                    embed.setColor(config.embedColor.r)
                    embed.setTitle('Premium Command')
                    embed.setDescription('The command you tried to use is only for premium servers. See the button below for more information.')

                    // Send the reply
                    await interaction.reply({ embeds: [embed], components: [row] })
                    // Return since the guild is not allowed to run the command due to a lack of premium
                    return
                }
            }

            obj.execute(client, interaction)
        }
    } else if (interaction.isMessageComponent() && interaction.componentType === 'BUTTON') {

        // Check if the used button is actually stored in the botCache object
        if (botCache.buttons.has(interaction.customId)) {

            // Retrieve the interaction data from the botCache object and run the binded function
            botCache.buttons.get(interaction.customId)(client, interaction)
        }
    }
}