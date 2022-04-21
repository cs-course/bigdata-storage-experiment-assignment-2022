#include <bits/stdc++.h>

using namespace std;

//0为正确性测试，有输出
//1为性能测试，无输出，有计时
#define TEST_MODE 1


class CuckooHash{
private:
    unsigned size;
    unsigned hashFuncNum;
    unsigned maxRelocateNum;
    double maxLoadFactor;
    char** hashTable;
    unsigned* cnt;
    list<char*> auxList;
    unordered_set<unsigned> posSet;
    unsigned totalRelocateNum=0;

    unsigned long long hashFunc1(const char* s)  {  
        unsigned long long hash = 1315423911;  
        for(int i = 0; s[i]!=0; i++)  {  
            hash ^= ((hash << 5) + s[i] + (hash >> 2));  
        }  
        return hash;  
    }
    unsigned long long hashFunc2(const char* s){
        int b     = 378551;  
        int a     = 63689;  
        unsigned long long hash = 0;  
        for(int i = 0; s[i]!=0; i++)  
        {  
            hash = hash * a + (int)s[i];
            a    = a * b;  
        }  
        return hash;  
    }  

    //随机迁移
    char* rndRelocate(char* s){
        unsigned long long h1=hashFunc1(s),h2=hashFunc2(s);
        unsigned long long h=h1;
        for(int i=0;i<this->hashFuncNum;i++,h+=h2){
            unsigned pos=h%this->size;
            if(this->hashTable[pos]==NULL){
                this->hashTable[pos]=s;
                return NULL; 
            }
        }
        h=h1;
        for(int i=0;i<this->hashFuncNum;i++,h+=h2){
            unsigned pos=h%this->size;
            if(posSet.find(pos)==posSet.end()){
                posSet.insert(pos);
                char* rnt=this->hashTable[pos];
                this->hashTable[pos]=s;
                totalRelocateNum++;
                return rnt;
            }
        }
        //所有位置都走过了，无法避免环路，返回原值
        return s;
    }

    typedef struct hrstStrct{
        char* s;
        unsigned cnt;
    }hrstStrct;

    //最少迁移次数启发式策略
    hrstStrct hrstRelocate_1(hrstStrct hs){
        unsigned long long h1=hashFunc1(hs.s),h2=hashFunc2(hs.s);
        unsigned long long h=h1;
        unsigned  minCnt=0x7FFFFFFF, minPos=-1;
        for(int i=0;i<this->hashFuncNum;i++,h+=h2){
            unsigned pos=h%this->size;
            if(this->hashTable[pos]==NULL){
                this->hashTable[pos]=hs.s;
                this->cnt[pos]=hs.cnt;
                return {NULL,0}; 
            }
        }
        h=h1;
        for(int i=0;i<this->hashFuncNum;i++,h+=h2){
            unsigned pos=h%this->size;
            if(cnt[pos]<minCnt && posSet.find(pos)==posSet.end()){
                minCnt=cnt[pos];
                minPos=pos;
            }
        }
        if(minPos!=-1){
            posSet.insert(minPos);
            hrstStrct rnt={this->hashTable[minPos],this->cnt[minPos]+1};
            this->hashTable[minPos]=hs.s;
            this->cnt[minPos]=hs.cnt;
            totalRelocateNum++;
            return rnt;
        }

        //所有位置都走过了，无法避免环路，返回原值
        return hs;
    }

    //最多空桶启发式策略
    char* hrstRelocate_2(char* s){
        unsigned long long h1=hashFunc1(s),h2=hashFunc2(s);
        unsigned long long h=h1;
        unsigned  minCnt=0xFFFF, minPos=-1;
        unsigned usedNum=0,emptyPos=-1;
        for(int i=0;i<this->hashFuncNum;i++,h+=h2){
            unsigned pos=h%this->size;
            if(this->hashTable[pos]!=NULL){
                usedNum++;
            }
            else if(emptyPos==-1){
                emptyPos=pos;
            }
        }
        if(emptyPos!=-1){
            this->hashTable[emptyPos]=s;
            this->cnt[emptyPos]=usedNum+1;
            return NULL; 
        }
        h=h1;
        for(int i=0;i<this->hashFuncNum;i++,h+=h2){
            unsigned pos=h%this->size;
            if(cnt[pos]<minCnt && posSet.find(pos)==posSet.end()){
                minCnt=cnt[pos];
                minPos=pos;
            }
        }
        if(minPos!=-1){
            posSet.insert(minPos);
            char* rnt=this->hashTable[minPos];
            this->hashTable[minPos]=s;
            this->cnt[minPos]=this->hashFuncNum;
            totalRelocateNum++;
            return rnt;
        }

        //所有位置都走过了，无法避免环路，返回原值
        return s;
    }


    void insertToAuxList(char* s){
        this->auxList.push_front(s);
        if(TEST_MODE==0)
            cout<<s<<" is relocated to linked list."<<endl;
        return;
    }


public:
    
    CuckooHash(unsigned size, unsigned hashFuncNum, unsigned maxRelocateNum,double maxLoadFactor){
        this->size=size;
        this->hashFuncNum=hashFuncNum;
        this->maxRelocateNum=maxRelocateNum;
        this->maxLoadFactor=maxLoadFactor;
        hashTable=(char**)malloc(sizeof(char*)*this->size);
        cnt=(unsigned*)malloc(sizeof(unsigned)*this->size);
        for(unsigned i=0;i<this->size;i++){
            hashTable[i]=NULL;
        }
        srand(time(NULL));
    }

    //返回该元素插入的位置在所有侯选位置中的序号
    //返回值为非负数表示插入成功，否则为插入失败
    int insert(const char* s){
        if(query(s)!=-1){
        //cout<<"The item has been inserted before."<<endl;
            return 2; //表示已经在表中
        }
        char* p=(char*)malloc(sizeof(char)*strlen(s+1));
        strcpy(p,s);
        unsigned reloceteNum=0;
        char* pre=NULL;
        do{
            pre=p;
            p=rndRelocate(p);
            //p不动，表示p的所有位置都已经走过，无法避免环路
            if(p==pre){
                break;
            }
            reloceteNum++;
        }while(p!=NULL && reloceteNum<this->maxRelocateNum);
        //cout<<s<<"is inserted to the hash table."<<endl;
        posSet.clear();
        if(p!=NULL){
            //表示无法解决哈希冲突，最后一个item插入辅助链表中
            insertToAuxList(p);
        }
        return 0;     //插入成功
    }

    //最少迁移次数启发式策略
    int hrstInsert_1(const char* s){
        if(query(s)!=-1){
            //cout<<"The item has been inserted before."<<endl;
            return 2; //表示已经在表中
        }
        char* p=(char*)malloc(sizeof(char)*strlen(s+1));
        strcpy(p,s);
        unsigned reloceteNum=0;
        char* pre=NULL;
        hrstStrct hs={p,0};
        do{
            pre=p;
            hs=hrstRelocate_1(hs);
            //p不动，表示p的所有位置都已经走过，无法避免环路
            if(hs.s==pre){
                break;
            }
            reloceteNum++;
        }while(hs.s!=NULL && reloceteNum<this->maxRelocateNum);
        //cout<<s<<"is inserted to the hash table."<<endl;
        posSet.clear();
        if(hs.s!=NULL){
            //表示无法解决哈希冲突，最后一个item插入辅助链表中
            insertToAuxList(hs.s);
        }
        return 0;     //插入成功
    }

    //最多空桶启发式迁移策略
    int hrstInsert_2(const char* s){
        if(query(s)!=-1){
            //cout<<"The item has been inserted before."<<endl;
            return 2; //表示已经在表中
        }
        char* p=(char*)malloc(sizeof(char)*strlen(s+1));
        strcpy(p,s);
        unsigned reloceteNum=0;
        char* pre=NULL;
        do{
            pre=p;
            p=hrstRelocate_2(p);
            //p不动，表示p的所有位置都已经走过，无法避免环路
            if(p==pre){
                break;
            }
            reloceteNum++;
        }while(p!=NULL && reloceteNum<this->maxRelocateNum);
        //cout<<s<<"is inserted to the hash table."<<endl;
        posSet.clear();
        if(p!=NULL){
            //表示无法解决哈希冲突，最后一个item插入辅助链表中
            insertToAuxList(p);
        }
        return 0;     //插入成功
    }

    //返回item在哈希表中的位置，若查找失败则返回-1
    int query(const char* s){
        unsigned long long h1=hashFunc1(s),h2=hashFunc2(s);
        unsigned long long h;
        for(int i=0;i<this->hashFuncNum;i++){
            if(i==0){
                h=h1;
            }
            else{
                h+=h2;
            }
            const char* p=this->hashTable[h%this->size];
            if(p!=NULL && strcmp(s,p)==0){
                //cout<<"Yes"<<endl;
                return h%this->size;
            }
        }
        for(auto iter=auxList.begin();iter!=auxList.end();iter++){
            if(strcmp(s,*iter)==0){
                return this->size;
            }
        }
        //cout<<"No"<<endl;
        return -1;
    }

    void clear(){
        for(unsigned i=0;i<this->size;i++){
            hashTable[i]=NULL;
        }
        memset(this->cnt,0,sizeof(this->cnt));
        auxList.clear();
        totalRelocateNum=0;
    }

    unsigned getTotalRelocateNum(){
        return this->totalRelocateNum;
    }

    unsigned getListLength(){
        return this->auxList.size();
    }

    void printConfig(){
        printf("{ size=%d  hashFuncNum=%d  maxRelocateNum=%d  load factor=%.2f }\n",this->size,this->hashFuncNum,this->maxRelocateNum,this->maxLoadFactor);
    }

    unsigned getSize(){
        return size;
    }

    double getMaxLoadFactor(){
        return maxLoadFactor;
    }

};

//正确性测试
void crtTest(CuckooHash& t,int modelTested, vector<string>& insertTestSet,vector<string>& queryTestSet){
    char filePath[13]={0};
    sprintf(filePath,"output_%d.txt",modelTested);
    FILE* fp=freopen(filePath,"w+",stdout);
    //cout<<insertTestSet.size()<<endl;
    int rnt=0;
    for(auto testItem : insertTestSet){
        switch (modelTested)
        {
        case 0:
            rnt=t.insert(testItem.c_str());
            break;
        case 1:
            rnt=t.hrstInsert_1(testItem.c_str());
            break;
        case 2:
            rnt=t.hrstInsert_2(testItem.c_str());
            break;
        default:
            break;
        }
        if(rnt==0){
            cout<<testItem<<" is inserted to the hash table."<<endl;
        }
        else if(rnt==2){
            cout<<testItem<< " has been inserted before."<<endl;
        }
    }    
    int i=0;
    for(auto testItem : queryTestSet){
        rnt=t.query(testItem.c_str());
        if(rnt>=0 && rnt<t.getSize()){
            cout<<testItem<<" in hash table."<<endl;
        }
        else if(rnt==t.getSize()){
            cout<<testItem<<" in auxiliary linked list."<<endl;
        }
        else{
            cout<<testItem<<" not found."<<endl;
        }
        //3个查询成功，1个查询失败，否则系统错误
        if((i%4==3&&rnt!=-1) || (i%4!=3&&rnt==-1)){
            if(rnt!=-1){
                cout<<"query result error!"<<endl;
                return;
            }
        }
        i++;
    }
    cout<<"0 error."<<endl;
    t.clear();
    fclose(fp);
}

//性能测试
void pfmcTest(CuckooHash& t, const int modelTested, vector<string>& insertTestSet){
    t.printConfig();
    time_t b=clock();
    switch (modelTested)
    {
    case 0:
        cout<<"Random selection strategy."<<endl;
        for(auto testItem: insertTestSet){
            t.insert(testItem.c_str());
        }
        break;
    case 1:
        cout<<"Minimum relocations selection strategy"<<endl;
        for(auto testItem: insertTestSet){
            t.hrstInsert_1(testItem.c_str());
        }
        break;
    case 2:
        cout<<"Maximum empty bucket selection strategy."<<endl;
        for(auto testItem: insertTestSet){
            t.hrstInsert_2(testItem.c_str());
        }
        break;
    default:
        break;
    }
    time_t e=clock();
    cout<<"Total run time: "<<e-b<<endl;
    cout<<"Total relocate num: "<<t.getTotalRelocateNum()<<endl;
    cout<<"Relocations per insert: "<<t.getTotalRelocateNum()/(t.getSize()*t.getMaxLoadFactor())<<endl;
    cout<<"Length of auxiliary linked list: "<<t.getListLength()<<endl<<endl;
    t.clear();
}


int main(){
    const unsigned size=1e7, hashFuncNum=6,maxRelocateNum=30;
    double maxLoadFactor=0.91;
    vector<string> insertTestSet, queryTestSet;
    int randNum=rand();
    char op;
    string s;
    int rnt;
    for(int i=0;i<size*maxLoadFactor;i++){
        s=to_string(randNum+i);
        insertTestSet.push_back(s);
        queryTestSet.push_back(s);
        if(i%3==2){
            s=to_string(randNum+i+size);
            queryTestSet.push_back(s);
        }
    }
    random_shuffle(insertTestSet.begin(),insertTestSet.end());
    CuckooHash t0(size,hashFuncNum,maxRelocateNum,maxLoadFactor);
    if(TEST_MODE==0){
        //第一次是让系统预热
        crtTest(t0,0,insertTestSet,queryTestSet);
        //第二次才正式开始测试
        crtTest(t0,0,insertTestSet,queryTestSet);
        crtTest(t0,1,insertTestSet,queryTestSet);
        crtTest(t0,2,insertTestSet,queryTestSet);
    }
    else if(TEST_MODE==1){
        //第一次是让系统预热
        pfmcTest(t0,0,insertTestSet);
        //第二次才正式开始测试
        pfmcTest(t0,0,insertTestSet);
        pfmcTest(t0,1,insertTestSet);
        pfmcTest(t0,2,insertTestSet);
    }
    return 0;
}
