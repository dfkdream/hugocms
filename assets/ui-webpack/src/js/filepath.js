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
        if (u1[u1.length-1]!=="/"){
            u1+="/";
        }
        if (u2[0]==="/"){
            u2=u2.slice(1)
        }
        return u1+u2;
    },

    abs: function(base, rel){
        let stack = base.split("/"),
            parts = rel.split("/");
        stack.pop(); // remove current file name (or empty string)
                     // (omit if "base" is the current folder without trailing slash)
        for (let i=0; i<parts.length; i++) {
            if (parts[i] === ".")
                continue;
            if (parts[i] === "..")
                stack.pop();
            else
                stack.push(parts[i]);
        }
        return stack.join("/");
    },

    clean: function(path){
        let stack = path.split("/");
        let result = [""];

        stack.forEach(i=>{
            switch(i){
                case "": case ".":
                    return;
                case "..":
                    result.pop();
                    break;
                default:
                    result.push(i);
            }
        });

        return this.join("/",result.join("/"));
    }
};
