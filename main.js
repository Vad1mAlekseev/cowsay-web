const { exec } = require('child_process');
const http = require('http');
fs = require('fs');
const { promisify } = require('util');

global._ = require('lodash');
const { APP_PORT, APP_HOST } = require('./app-config')
const readFile = promisify(fs.readFile);

const server = http.createServer(async (req, res) => {
    res.statusCode = 200;

    if (!req.url.slice(1)) {
        let template = (await readFile('templates/index.html')).toString();
        res.end(await new Promise((resolve, reject) => {
            exec(`cowsay -l`, (error, stdout, stderr) => {
                if (stderr || error) {
                    res.end(`exec error: ${error}`);
                    resolve(stderr);
                    re2turn;
                }
                const buttons =
                    _(stdout)
                        .split('\n')
                        .splice(1)
                        .words()
                        .value();
                resolve(_.template(template)({ buttons }));
            });
        }));
        return;
    }

    res.end(await new Promise((resolve, reject) => {
        exec(`cowsay -f ${req.url.split('').splice(1).join('')} hello`, (error, stdout, stderr) => {
            if (stderr || error) {
                res.end(error.toString());
                resolve(stderr);
                return;
            }
            resolve(stdout);
        });
    }));
});

server.listen(APP_PORT, APP_HOST, () => {
    console.log(`Server listen ${APP_HOST}:${APP_PORT}`);
});
