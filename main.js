const Koa = require('koa');

global._ = require('lodash');
const { APP_PORT, APP_HOST } = require('./app-config');
const router = require('./router');
const app = new Koa();

async function start() {
  app.use(router.routes());
  app.listen(APP_PORT, APP_HOST, () => {
    console.log(`Server listen ${APP_HOST}:${APP_PORT}`);
  });
}

start();
