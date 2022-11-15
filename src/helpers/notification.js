const OneSignal = require("@onesignal/node-onesignal");

const user_key_provider = {
  getToken() {
    return process.env.ONESIGNAL_USER_KEY;
  },
};

const app_key_provider = {
  getToken() {
    return process.env.ONESIGNAL_REST_KEY;
  },
};

const configuration = OneSignal.createConfiguration({
  authMethods: {
    user_key: {
      tokenProvider: user_key_provider,
    },
    app_key: {
      tokenProvider: app_key_provider,
    },
  },
});

const client = new OneSignal.DefaultApi(configuration);

const addPlayer = async ({ email, onesignal_id }) => {
  const player = new OneSignal.Player();
  player.device_type = 1;
  player.app_id = process.env.ONESIGNAL_APP_ID;
  player.identifier = onesignal_id;
  player.external_user_id = email;
  await client.createPlayer(player);
};

const deletePlayer = async ({ onesignal_id }) => {
  await client.deletePlayer(process.env.ONESIGNAL_APP_ID, onesignal_id);
};

const sendNotification = async ({ title, subtitle, playerId }) => {
  const notification = new OneSignal.Notification();

  // App ID
  notification.app_id = process.env.ONESIGNAL_APP_ID;

  // Icon App
  notification.small_icon = "ic_launcher";

  // Title
  notification.headings = {
    en: title,
  };

  // Subtitle
  notification.contents = {
    en: subtitle,
  };

  // Target
  notification.include_player_ids = [playerId];

  await client.createNotification(notification);
};

module.exports = { sendNotification, addPlayer, deletePlayer };
