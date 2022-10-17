const fs = require('fs');

module.exports = {
    deleteFile: (path) => {
        if (path != "" && path != null) {
            fs.unlink(path, (err) => {
                if (err) throw err
            });
        };
    },
}