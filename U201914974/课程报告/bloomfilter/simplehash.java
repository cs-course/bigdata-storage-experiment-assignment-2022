package bloomfilter;

public class simplehash {
    private int seed;
    private int m;
    simplehash(int m,int seed){
        this.m = m;
        this.seed = seed;
    }
    public int hash(String data){
        int code = 0;
        for(int i = 0; i < data.length(); i++){
            code = code * this.seed + (int)data.charAt(i);
        }
        return Math.abs(code) % m;
    }
}
