# Data
Information about all the data that is being collected is displayed below. Guild related data will be removed if the bot gets kicked out of the guild.

### Guilds
guild_id (string):
 - staff_role (string)
 - sug_channel (string)
 - rep_channel (string)
 - upvotes (user id array)
 - downvotes (user id array)

### Notes
note_id (string):
 - author_hash (string)
 - title (string)
 - contents (string)
 - timestamp (int64)

### Interactive
id (string):
 - expiry (timestamp)
 - guild_id (string)
 - channel_id (string)
 - message_id (string)
 - participants (string[])

### Caching
Cached items will be deleted after 8 hours.

Suggestions: id -> upvotes, downvotes
Interactive: id -> expiry, participants