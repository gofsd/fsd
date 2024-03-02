const htmlparser2 = require("htmlparser2");
const fs = require("fs")
var word = {};
var wordList = [];
var counter = 0;
var freqList;
var coll;

const parser = new htmlparser2.Parser({
    onopentag(name, attributes) {
        /*
         * This fires when a new tag is opened.
         *
         * If you don't need an aggregated `attributes` object,
         * have a look at the `onopentagname` and `onattribute` events.
         */


    },
    ontext(text) {
        /*
         * Fires whenever a section of text was processed.
         *
         * Note that this can fire at any point within text and you might
         * have to stich together multiple pieces.
         */
        if (text.charCodeAt() == 13) {
         return;
        }
        if(text == ", ") {
        return;
        }



        if(["n.", "adj.", "exclam.", "adv.", "prep.", "v.", "pron.", "conj."].includes(text)){
            if(word.partsOfSpeech == undefined) {
                word.partsOfSpeech = []
            }
            if(text.includes("n.") && !text.includes( "pron.")) {
                word.partsOfSpeech.push("noun");
            }
            if(text.includes("adj.")) {
                word.partsOfSpeech.push("adjective");
            }
            if(text.includes("exclam.")) {
                word.partsOfSpeech.push("interjection");
            }
            if(text.includes("adv.")) {
                word.partsOfSpeech.push("adverb");
            }
            if(text.includes("prep.")) {
                word.partsOfSpeech.push("preposition");
            }
            if(text.includes("v.") && !text.includes("adv.")) {
                word.partsOfSpeech.push("verb");
            }
            if(text.includes( "pron.")) {
                word.partsOfSpeech.push("pronoun");
            }
            if(text.trim().includes("conj.")) {

                word.partsOfSpeech.push("conjunction");
            }

        }

        if(["A1", "A2", "B1", "B2", "C1", "C2"].includes(text.trim())) {
            counter++;
            word.level = text.trim();
            if(word.partsOfSpeech && word.partsOfSpeech.length){
                                wordList.push(word);

            }

            word = {};
        }
        if(word.name == undefined) {
            word.name = !text.includes("1 ") ? text.trim() : text.replace("1 ", "");
        }
    },
    onclosetag(tagname) {
        /*
         * Fires when a tag is closed.
         *
         * You can rely on this event only firing when you have received an
         * equivalent opening tag before. Closing tags without corresponding
         * opening tags will be ignored.
         */
        if (tagname === "div") {
            word = {};
        }
        if(tagname === "body") {
            wordList = wordList.map(word => {
                word.frequency = freqList.findIndex((str) => str.trim()==word.name.trim());
                return word;
            }).reduce((ac, it)=>{
                if(it.frequency > 0 && it.frequency < 20000){
                    it.examples = [];
                    it.transcription = "";
                    it.meanings = [];
                    ac.push(it)
                }
                return ac;
            }, []).map((item) => {
                  var colItem = null;
                  coll.forEach((str) => {
                      obj = str == "" ? null : JSON.parse(str);
                      if(obj && obj.word.toLowerCase() == item.name.toLowerCase()) {
                          colItem = obj;
                      }
                  });
                  if(colItem && colItem.example_word && colItem.meaning_word){
                      item.examples.push(colItem.example_word.replace("<i>", "").replace("</i>", "").replace("<b>", "").replace("</b>", ""));
                      item.meanings.push(colItem.meaning_word.replace("<i>", "").replace("</i>", "").replace("<b>", "").replace("</b>", ""));
                      item.transcription = colItem.transcription_word;
                  }
                  return item;
              });
            wordList.sort((a,b)=> {
                return a.frequency > b.frequency ? -1 : 1;
            })
            console.log(wordList.length)
            fs.writeFile("./output_data/test.json", JSON.stringify(wordList), "utf8", (err)=> console.log(err));
        }
    },
});

fs.readFile('./input_data/google_freak_20k.txt', 'utf8',function (err,str) {
    fs.readFile('./input_data/wordColection.json', 'utf8',function (err,strCol) {
        freqList = str.split("\n");
        coll = strCol.split("\n")
        fs.readFile('./input_data/ox3000.html', 'utf8', function (err,data) {
          if (err) {
            return console.log(err);
          }
          parser.write(
              data
          );
        });

        fs.readFile('./input_data/ox5000.html', 'utf8', function (err,data) {
          if (err) {
            return console.log(err);
          }
          parser.write(
              data
          );
        });
    });
})

parser.end();