package bloomfilter;

import java.util.BitSet;

public class bloomfilter {
    private int m;
    private int[] seeds;
    private simplehash[] hashs;
    private BitSet bits;

    bloomfilter(int m,int[] seeds){
        this.m = m;
        this.bits = new BitSet(m);
        this.seeds = seeds;
        this.hashs = new simplehash[this.seeds.length];
        for(int i = 0; i < this.seeds.length; i++){
            this.hashs[i] = new simplehash(this.m,this.seeds[i]);
        }
    }

    public void add(String data){
        for(simplehash h : this.hashs){
            this.bits.set(h.hash(data),true);
        }
    }

    public boolean contains(String data){
        for(simplehash h : this.hashs){
            if(!this.bits.get(h.hash(data))) return false;
        }
        return true;
    }
}
