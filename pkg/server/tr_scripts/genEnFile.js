var unirest = require("unirest");
var _ = require('lodash');

var fs = require("fs");

fs.readFile('./input_data/lookup/uk.json', 'utf8', async function (err,str) {
    ukMap = JSON.parse(str);
    ukMap = ukMap.map(function(item, idx, arr) {
        newItem = {};
        newItem.normalizedSource = item.normalizedSource;
        newItem.displayTarget = item.displaySource;
        newItem.translations = item.translations.reduce(function(ac, it){
            it.backTranslations.forEach(function(it){
                if (item.normalizedSource != it.normalizedText) {
                    it.displayTarget = it.normalizedText
                    ac.push(it);
                }
            });
            return ac;
        }, [])
        return newItem;
    });
    fs.writeFile(`./input_data/lookup/en.json`, JSON.stringify(ukMap), "utf8", (err)=> console.log(err));

});