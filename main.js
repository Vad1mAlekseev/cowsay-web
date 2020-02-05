const childProcess = require('child_process');
const Koa = require('koa');
fs = require('fs');
const { promisify } = require('util');
const Router = require('koa-router');

global._ = require('lodash');
const { APP_PORT, APP_HOST } = require('./app-config');
const readFile = promisify(fs.readFile);
const exec = promisify(childProcess.exec);
const app = new Koa();
const router = new Router();

async function start() {
  const sowsay = async (flags, ...args) => {
    const { stdout, stderr } = await exec(`cowsay -${flags} ${args.join(' ')}`);
    return stderr || stdout;
  };
  const templates = {
    index: await readFile('templates/index.html'),
    cow: await readFile('templates/cow.html'),
  };

  router.get('/', async ctx => {
    try {
      const cowsayResult = await sowsay('l');

      const cows = _(cowsayResult)
        .split('\n')
        .splice(1)
        .join(' ')
        .trim()
        .split(' ');

      ctx.body = _.template(templates.index)({ cows });
    } catch (err) {
      ctx.throw(400, err);
    }
  });

  router.get('/:file', async ctx => {
    const cow = await sowsay(
      'f',
      ctx.params.file,
      ctx.query.text || 'Hello World!'
    );
    ctx.body = _.template(templates.cow)({ cow });
  });

  app.use(router.routes());
  app.listen(APP_PORT, APP_HOST, () => {
    console.log(`Server listen ${APP_HOST}:${APP_PORT}`);
  });
}

start();
