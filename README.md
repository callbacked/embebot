# Embe Bot
[![](https://github.com/callbacked/embebot/blob/main/assets/add.png)](https://discord.com/api/oauth2/authorize?client_id=1100908930458198098&permissions=274877983744&scope=applications.commands%20bot)

A quick and silly little Python bot I created that utilizes services such as vxtwitter, vxtiktok, vxreddit, and ddinstagram in one bot to properly embed media for the user.

I probably won't update it as frequently since its in a good enough state for me, but perhaps I can improve 
upon it over time.

# How it works

Rather than re-inventing the wheel and creating a way to fix embeds on Discord for every major social media site -- I opted to automate the manual way of using these aforementioned services so it is done for you.

![](https://github.com/callbacked/embebot/blob/main/assets/manual-embed.gif)

## Simply paste your link and send.
![](https://github.com/callbacked/embebot/blob/main/assets/embed.gif)

**Just post a link from either twitter.com, x.com, tiktok.com, instagram.com (reels only for now), and Embe Bot will detect it and reply to the original message with the converted link.**


## Building and hosting it yourself (On Docker)

**This will require you to make an application in the [Discord Developer Portal](https://discord.com/developers/applications)**

	- Create a New Application > Go to the Bot tab
	- Turn on all Privileged Gateway intents
	- Under Bot Permissions, give it Text Permissions
	- Add bot to your server 
- Clone the github repository
- in the /embebot directory build the image via ``docker build .``
- in docker-compose.yml modify the line ``DISCORD_BOT_TOKEN =`` and add in your own Bot token by pressing "Reset Token" in the Bot tab and copy it to the line.

## Want to use your own endpoints? No worries

Whether it'd be vx/fxtwitter, ddinstagram, vxreddit, sometimes these services get rate limited which can lead to a much slower experience. Because of this you have the option to change the embed link to something else to avoid congestion, which is handy for those who already self-host these services.

Just edit the ``config.ini`` file, set ``EndpointOverride`` to ``true`` and start adding in your own endpoints for any of the services listed. Should you want to switch back, set them as default.

For reference, the default embed endpoints are:

* Twitter/X - vxtwitter.com
* Instagram - ddinstagram.com
* Reddit - vxreddit.com





## Thanks
This bot definitely wouldn't even exist without the heavy lifting of these projects, I just merely combined them to make the Discord experience slightly more bearable. So I thank these people orders of magnitude smarter than me for their great work:

* [vxtiktok - dylanpdx](https://github.com/dylanpdx/vxtiktok)
* [vxreddit - dylanpdx](https://github.com/dylanpdx/vxReddit)
* [BetterTwitFix - dylanpdx](https://github.com/dylanpdx/BetterTwitFix)
* [InstaFix - Wikidepia](https://github.com/Wikidepia/InstaFix)
# 
Note: this thing is licensed under do whatever you want I dont care



