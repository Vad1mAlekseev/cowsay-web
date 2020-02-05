const Router = require('koa-router');
const childProcess = require('child_process');
const fs = require('fs');
const { promisify } = require('util');

const exec = promisify(childProcess.exec);
const router = new Router();

const templates = {
  index: fs.readFileSync('templates/index.html'),
  cow: fs.readFileSync('templates/cow.html'),
};
const static = {
  ico: fs.readFileSync('static/favicon.ico'),
};

async function cowsay(flags, ...args) {
  const { stdout, stderr } = await exec(`cowsay -${flags} ${args.join(' ')}`);
  return stderr || stdout;
}

async function getCows() {
  const cowsayResult = await cowsay('l');
  return _(cowsayResult)
    .split('\n')
    .splice(1)
    .join(' ')
    .trim()
    .split(' ');
}

function defineRoutes() {
  router.get('/favicon.ico', async ctx => {
    ctx.type = 'image/x-icon';
    ctx.body = static.ico;
  });

  router.get('/', async ctx => {
    try {
      ctx.body = _.template(templates.index)({ cows: await getCows() });
    } catch (err) {
      ctx.throw(400, err);
    }
  });

  router.get('/:cow', async ctx => {
    const cows = await getCows();
    const requestedCow = ctx.params.cow;
    const nextCowIdx = _.indexOf(cows, requestedCow) + 1;
    const prevCowIdx = _.indexOf(cows, requestedCow) - 1;

    const cow = await cowsay(
      'f',
      ctx.params.cow,
      requestedCow ? `"${requestedCow}"` : `"I'm a ${requestedCow}!"`
    );
    ctx.body = _.template(templates.cow)({
      cow,
      nextCow: cows[nextCowIdx],
      prevCow: cows[prevCowIdx],
    });
  });

  return router;
}

module.exports = defineRoutes();
