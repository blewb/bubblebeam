I have to log my hours for my job in a system called Streamtime. I currently handle this by tracking my hours through a simple text format, than manually enter them at the end of the day/week. I have been working on this little Golang TUI to make this process much more efficient. It will parse my text file, then open an interface where I can quickly select the correct "jobs" that my time should map to.

I created the <span> module to do the parsing. I believe this is working as intended, with a few simple tests.

I have started writing the <stream> module to be the API layer. Some basic elements are there but it is incomplete.

I have half built the TUI part inside /cmd/app. I would like you to build a new TUI inside /cmd/app2. Use Bubbletea and Lipgloss for the styling. It needs to be responsive but I am OK with setting a minimum size for the window width, below which the app does not render and just shows a warning to increase the window size.

There are multiple steps/layers to walk through to get this working.

1. Boot the TUI model, and parse the hours logging info from a given text file (currently there are a few in /sample)
2. Load the list of active jobs from the streamtime API - for ease of development I have saved a sample response as a JSON file in /temp so we don't have to pull from the API each time, during development.
3. The hours being logged might be for a single day, or for multiple days. Need an interface to select which date(s) are being referenced. Always assume the day name refers to the most recent date of that name (including the current). Can also assume that if, for example, Monday is locked in for a date, then Tuesday would mean the next date, Wednesday the date after. However still need an option to edit them individually.
4. The main body of the interface will be split into 3 panels - the largest will fill the top half of the screen, full width, and contain a table with the hours entries parsed from the text file, for the selected day. This panel will be in focus initially, with the user able to navigate up and down using arrow keys. Left/right arrows will change the day, if there are multiple days.
  a. Panel 2 will sit in the bottom half, and use about 2/3rds of the width - once an hour entry from panel 1 is selected this panel takes focus and lets the user do a case insensitive search of the active jobs to select the correct one.
  b. Panel 3 will also sit in the bottom half and fill the remaining 1/3rd of the width. When an active job is selected this panel will take focus and will need to fetch the available job items for the chosen job, then let the user confirm which is the correct job item.
5. Once a job item is confirmed we focus back to Panel 1 and continue the process (in later stages we'll probably refine this to automatically jump to panel 2 with the next hour log selected).
6. The job item selections made from panel 3 will need to be saved into some data structure that maps to the hours items, and then later we'll work through actually submitting those entries to the streamtime API.

For the 3 panels I want them each to have a small header with room for a title (for now just Panel 1, Panel 2, Panel 3) in the top left, and an icon in the top right which can indicate changes of status (such as a loader when fetching data). The panels will have a mid-grey border in default state, changing to a solid white border when focused.

I am open to suggestions for sensible key bindings for interacting with this TUI. E.g. using tab to navigate between panels?

You can refer to my original /cmd/app version for some sense of where I was heading, although it is incomplete. The app's name is BUBBLEBEAM, and I had a colourful title in my original version which you can copy and sit above all the panels.
