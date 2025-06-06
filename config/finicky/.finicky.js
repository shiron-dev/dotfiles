module.exports = {
  defaultBrowser: "Arc",
  handlers: [
    {
      match: ["discord.com*"],
      url: {
        protocol: "discord",
      },
      browser: "Discord",
    },
    {
      match: finicky.matchHostnames([
        "drive.google.com",
        "docs.google.com",
        "forms.gle"
      ]),
      browser: {
        name: "Google Chrome",
        profile: "Profile 1",
      },
    },
  ],
};
