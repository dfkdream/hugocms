/*
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
module.exports = {
    ext: function(path){
        return path.split(".").pop()
    },

    join: function(u1, u2){
        if (u1[u1.length-1]!="/"){
            u1+="/";
        }
        if (u2[0]=="/"){
            u2=u2.slice(1)
        }
        return u1+u2;
    }
};