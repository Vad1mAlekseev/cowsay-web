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
const staticFiles = {
  ico: fs.readFileSync('static/favicon.ico'),
};

async function cowsay(flags, text) {
  let query = `cowsay ${_(flags).map((flag) => `-${flag}`).join(' ')}`;
  if (text) {
    if (text === 'random') {
      query = `fortune | ${query}`;
    } else {
      query = `echo ${_(text).words().join(' ')} | ${query}`;
    }
  }
  const { stdout, stderr } = await exec(query);
  return stderr || stdout;
}

async function getCows() {
  const cowsayResult = await cowsay(['l']);
  return _(cowsayResult)
    .split('\n')
    .splice(1)
    .join(' ')
    .trim()
    .split(' ');
}

function defineRoutes() {
  router.get('/favicon.ico', async (ctx) => {
    ctx.type = 'image/x-icon';
    ctx.body = staticFiles.ico;
  });

  router.get('/', async (ctx) => {
    try {
      ctx.body = _.template(templates.index)({ cows: await getCows() });
    } catch (err) {
      ctx.throw(400, err);
    }
  });

  router.get('/:cow', async (ctx) => {
    const cows = await getCows();
    const requestedCow = ctx.params.cow;

    if (!_.includes(cows, requestedCow)) {
      return;
    }

    const cow = await cowsay([`f ${requestedCow}`], ctx.query.text || 'random');
    if (ctx.query.mode === 'simple') {
      ctx.body = cow;
      return;
    }
    const nextCowIdx = _.indexOf(cows, requestedCow) + 1;
    const prevCowIdx = _.indexOf(cows, requestedCow) - 1;

    ctx.body = _.template(templates.cow)({
      cow,
      nextCow: cows[nextCowIdx],
      prevCow: cows[prevCowIdx],
    });
  });


  return router;
}

module.exports = defineRoutes;
