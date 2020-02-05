const APP_HOST = process.env.APP_HOST || '0.0.0.0';
const APP_PORT = process.env.APP_PORT || process.env.PORT || '4040';

module.exports = {
    APP_HOST, APP_PORT,
}
