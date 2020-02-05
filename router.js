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

function defineRoutes() {
  const cowsay = async (flags, ...args) => {
    const { stdout, stderr } = await exec(`cowsay -${flags} ${args.join(' ')}`);
    return stderr || stdout;
  };

  router.get('/favicon.ico', async ctx => {
    ctx.type = 'image/x-icon';
    ctx.body = static.ico;
  });

  router.get('/', async ctx => {
    try {
      const cowsayResult = await cowsay('l');

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

  router.get('/:cow', async ctx => {
    const cow = await cowsay(
      'f',
      ctx.params.cow,
      ctx.query.text || 'Hello World!'
    );
    ctx.body = _.template(templates.cow)({ cow });
  });

  return router;
}

module.exports = defineRoutes();
