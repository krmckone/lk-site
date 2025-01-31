# A Short Perspective on My Time with Steam Deck LCD and OLED

This is a short essay on my time so far with the Steam Deck, both LCD and OLED versions.

### Receiving the Steam Deck LCD

I first received my Steam Deck LCD 256GB in July of 2022. I had been waiting for about 4 or 5 months at this point after putting five dollars down to reserve one.

![account info](../images/steam_deck/account_info.png "Account Listing")

![shipping info](../images/steam_deck/shipping_info.png "Shipping Info")

In summer of 2022 I was excited to get the steam deck for a few reasons. Since the Deck's announcement the year prior, I figured that having most of my Steam library available on a handheld would be a huge game changer and I'll get into why that is later. In addition to games, it was just a little linux machine that made it easy to get to the regular linux desktop. It has a separate dock that you can plug your USB devices and monitor into. Even if I'm not one to hack away at linux all day, I think it's really good that a reliable and easy to use Linux device is now so easy get your hands on. On the controls front, the track pads were an improvement over the steam controller's that came out around 2015 and the back paddles added additional function while letting you keep your thumbs on the track pads or sticks. The LCD screen at the time was honestly not something I considered either good or bad, however, the later OLED version that we'll get into is a clear improvement over the first iteration's. I'll mention this more than once, but I work at the computer for multiple hours at a time during the day. Sometimes the last thing I want to do at the end of a work day is continue to sit at my desk and play PC games. Even though I really enjoy gaming on my desktop with a 7900 XT and 170Hz monitor, one can only ensure so many hours looking at the same screen in the same room in the same chair per day.

### What Kind of Games Work Well

Games with native linux support work great with no tweaking necessary in general. Proton also enables tons of windows-built games to run just fine on Steam Deck. Some games over time lose proton compatibility through later updates but in general proton is a great bit of software in the linux gaming space. On that front though Valve has put in tons of work to make all of the software around the controller configuration system powerful and exposes a ridiculous amount of options, particularly for the track pads and joy sticks. When I wanted to try a game without native controller support, I did have to take 10-20 minutes to mess with some options. When playing FPSs, I want to use the right track pad as either mouse or joystick-like-mouse. The sensitivity is important to get comfortable although I still don't have the perfect formula for this. I've been finding that most games are a little different in how they handle mouse input and things may get lost in translation through all of the software layers handling the steam deck input. Regardless, it's pretty easy to get a right track pad configured per game that works well. A lot of users also recommend a track pad + gyro setup for FPSs although I still have not really gotten the hang of that.

![steam deck lcd](../images/steam_deck/steam_deck_lcd.jpg "Steam Deck LCD")

There's lots of open source plugins as well that you can get for your deck, like one that shows you information from the third party ProtonDB for each of your games. Sometimes user-provided data on ProtonDB is more accurate than Valve's internal testing, so this is useful to have handy. The steam deck's native resolution is 1280\*800 at 16:10 so games that natively support this resolution will not cause any black bars around the deck's screen. As an alternative, you can set the in-game resolution to 1280x720 and use the deck's built in scaling mode to remove the black bars. This comes at the cost of stretching/distoring the image, though.

Here is my summary/generalization for features that should allow the game to play well on Steam Deck:

- Highly configurable graphics
  - By default the game should pick settings that allow at minimum 30 FPS on average
  - If the game does not default to good performance settings, it needs to allow you to tweak them
- Supports native 16:10 resolution
  - This avoids margins/black bars around the screen
- Has either native Linux support or runs normally under Proton with minimal performance reduction
  - Some games may require you to manually pick the Proton version in Steam before it will boot
- Has native controller support
  - This honestly is not a deal-breaker for me, but having native controller support is a huge QoL feature
  - If the game does not have controller support, it should allow key remapping to facilitate a custom community controller configuration
- For multiplayer games, supports a Linux version of any anti-cheat
  - Games that only run under a native Windows environment due to their anti-cheat will not run as expected on Steam Deck

One additional feature of the steam deck that may go overlooked is the support for multiple controllers simultaneously, just like a console. For example, split screen in local-coop games should just work if you connect multiple controllers to the deck. This works great with the steam deck dock plugged into a TV.

## My top 50 most played games on steam deck

#### This list is automatically updated twice daily

<div>
  {{ template "steamDeckTop50" . }}
</div>

### What Kind of Games Do Not Work Well

Not all of the bullet points above should be given the same weight, but games that I enjoy the most on Steam Deck basically meet all of those. Occasionally I will play a game that maybe only supports 16:9 resolutions but that is probably the most minimal of compromises to make. If your game requires a ton of controller configuration mapping or runs poorly performance-wise, that's something that is better played on your more powerful gaming PC with mouse and keyboard.

### Things about the LCD that were not my favorite

- It's a little heavy
  - During play sessions of over thirty minutes, my arms and wrists would get tired unless I had something to rest the Steam Deck on
- Battery life is not great but this also depends on game/graphics settings. There are also a lot of power-saving options in the three-dot menu
- The LCD screen is not as good as the OLED version with worse colors/brightness and refresh rate. However, at LCD launch this was not a big deal with me
- My left track pad was clicky but I think I just had a less-than-ideal quality of build. My guess is Valve's manufacturing process has improved in the lifetime of the LCD
- My left trigger was squishy. I ended up taking the back plate off and using one of my wife's nail filers to reduce the height of the plastic around the left trigger housing. This fixed the squishy issue for me

### Why did I get an OLED

![OLED shipping info](../images/steam_deck/oled_shipping_info.png "OLED Shipping Info")

For Christmas 2023 I got my wife an OLED version so we could play co-op games together on the couch. Since I was too excited to let my wife open it up herself, I broke into the package in early December to turn it on for the first time. I was surprised with the difference in weight and improvement in screen. However, I was still on the fence about upgrading my LCD. It only took 6 months for my jealousy to boil over after watching my wife play her OLED while I still used my LCD. In May 2024 I sold my LCD on Facebook marketplace for $260 and ordered an OLED 512GB.

Basically, jealousy.

### Should you get an OLED Steam Deck or LCD?

Some people online have been saying the OLED is at most a side-grade to the LCD.

In my opinion the OLED is an upgrade in its own right. I am very pleased with the perceived decrease in weight on the OLED and I notice the difference over longer play sessions. If you already have an LCD and enjoy it then I would say stick with it, but I personally wanted the better screen and reduced weight. The improved Wifi and battery are also nice additions.

- The screen is also a big upgrade in itself
  - Colors are much nice and I can perceive the difference in brightness settings much better compared to the LCD
  - I don't own a Switch OLED but I imagine the difference in those units' screens are similar enough
- I find the physical feedback on the track pads much better on the OLED
  - At the time of writing I'm not sure if they used different track pads for OLED but compared to my old LCD I enjoy using them much better.
- Overall I also just like the color scheme better on the OLED
  - The all-black joysticks are sharp and the subtle orange power button is also sick.

![steam deck OLED](../images/steam_deck/oled.jpg "Steam Deck OLED")

I was able to list my dislikes about the LCD around the screen and weight because the OLED was able to bring obvious improvements to those areas. Without the OLED the comparison would not be possible.
