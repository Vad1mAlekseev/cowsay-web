const { exec } = require('child_process');
const Koa = require('koa');
const app = new Koa();
fs = require('fs');
const { promisify } = require('util');

global._ = require('lodash');
const { APP_PORT, APP_HOST } = require('./app-config');
const readFile = promisify(fs.readFile);

app.use(async ctx => {
  const URL = ctx.request.url;

  if (!URL.slice(1)) {
    let template = await readFile('templates/index.html');
    ctx.body = await new Promise((resolve, reject) => {
      exec(`cowsay -l`, (error, stdout, stderr) => {
        if (stderr || error) {
          ctx.body = `exec error: ${error}`;
          resolve(stderr);
          re2turn;
        }
        const buttons = _(stdout)
          .split('\n')
          .splice(1)
          .words()
          .value();
        resolve(_.template(template)({ buttons }));
      });
    });
    return;
  }

  ctx.body = await new Promise((resolve, reject) => {
    exec(
      `cowsay -f ${URL.split('')
        .splice(1)
        .join('')} hello`,
      (error, stdout, stderr) => {
        if (stderr || error) {
          ctx.body = error.toString();
          resolve(stderr);
          return;
        }
        resolve(stdout);
      }
    );
  });
});

app.listen(APP_PORT, APP_HOST, () => {
  console.log(`Server listen ${APP_HOST}:${APP_PORT}`);
});
