const BloomFilter = require('./BloomFilter')


class PBF_HT {
  SBFArr; // the SBF of each property
  HT; // hash table
  p;  // the number of properties
  n;  // the max capacity of each SBF
  error;  // the error rate
  constructor(p, n, error){
    this.p = p;
    this.n = n;
    this.error = error;
    this.initSBFArr();
    this.initHT();
  }

  initSBFArr(){
    this.SBFArr = new Array(this.p);
    for (let i = 0; i < this.p; i++) {
      this.SBFArr[i] = new BloomFilter(this.n, this.error);
    }
  }

  initHT(){
    this.HT = new Map();
  }

  add(el){
    let v = 0;
    for (let i = 0; i < this.p; i++) {
      const prop = el[i];
      this.SBFArr[i].add(prop)
      v += this.SBFArr[i].getVerificationValue(prop);
    }
    this.HT.set(v, true);
  }

  contain(el) {
    let v = 0;
    for (let i = 0; i < this.p; i++) {
      const prop = el[i];
      if (!this.SBFArr[i].contain(prop)) {
        return false;
      }
      v += this.SBFArr[i].getVerificationValue(prop);
    }
    if (this.HT.has(v)) {
      return true
    } else {
      return false;
    }
  }

  getHT(){
    return this.HT;
  }
}

const PBF_HT_INSTANCE = new PBF_HT(2, 10, 0.01);

PBF_HT_INSTANCE.add(['large','red'])
PBF_HT_INSTANCE.add(['small','green'])
console.log(PBF_HT_INSTANCE.contain(['large','red']));
console.log(PBF_HT_INSTANCE.contain(['small','green']));
console.log(PBF_HT_INSTANCE.contain(['large','green']));
console.log(PBF_HT_INSTANCE.contain(['small','red']));
console.log(PBF_HT_INSTANCE.getHT());