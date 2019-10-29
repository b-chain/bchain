#pragma once
#ifdef __cplusplus
extern "C"
{
#endif
#define OUT
    // console api
    int log(char *msg);

    // auth api
    bool isHexAddress(char *address);
    bool requireAuth(char *address);
    bool isAccount(char *address);
    bool isContract(char *address);

    // system api
    int hex2Int(char *hex);

    // assert api
    void assert(bool test, char *assertMsg);

    // crypo api
    void sha1(char *in, int len, OUT char *hash);
    void sha256(char *in, int len, OUT char *hash);
    void sha512(char *in, int len, OUT char *hash);
    void recover(char *hexMsg, char *hexSig, OUT char *address);

    // producer api
    void block_producer(OUT char *number);
    long long block_number(void);
    void requireRewordAuth(void);

    // action api
    void action_sender(OUT char *sender);
    void contract_address(OUT char *self);
    void contract_creator(OUT char *creator);

    // contract api
    void contract_create(char *creator, char *code, OUT char *addr);

    typedef struct stTopic
    {
        char *data;
        int data_len;
        stTopic *next;
    } stTopic;

    typedef struct stData
    {
        char *data;
        int data_len;
        stData *next;
    } stData;

    void emitEvent(stTopic *topics, stData *datas);

    // memory api
    typedef unsigned int size_t;
    void *memset(void *s, int ch, size_t n);
    void *memcpy(void *dest, const void *src, size_t n);
    void *memmove(void *dest, const void *src, size_t count);
    int memcmp(const void *buf1, const void *buf2, unsigned int count);

    // database api
    void db_emplace(char *key, int key_len, char *val, int val_len);
    void db_set(char *key, int key_len, char *val, int val_len);
    int db_get(char *key, int key_len, char *val);

    // cache database api
    void cacheDb_emplace(char *key, int key_len, char *val, int val_len);
    int cacheDb_get(char *key, int key_len, char *val);

    // result api
    void setResult(char *date, int data_len);

    // call api
    void action_call(char *addr, char *para);
    void contract_call(char *addr, char *para);

#define TypeI32 0
#define TypeI64 1
#define TypeF32 2
#define TypeF64 3
#define TypeAddress 4

#define SysAddressLen 43
    typedef struct stCallPara
    {
        int type; //TypeI32 ... TypeAddress
        void *data;
        int data_len;
        stCallPara *next;
    } stCallPara;
    void getCallPara(char *funcName, stCallPara *para, OUT char *out);

    void action_callWithPara(char *addr, char *funcName, stCallPara *para);
    void contract_callWithPara(char *addr, char *funcName, stCallPara *para);

    //string api
    int strlen(char *in);
    int strjoint(char *firstStr, char *secondStr, OUT char *out);
    void str2lower(char *str);
    int strcmp(char *firstStr, char *secondStr);

    //weight api
    unsigned long long getWeight(unsigned long long);

    //bigint api
    int big_add(char *a, char *b, char *ret);
    int big_sub_safe(char *a, char *b, char *ret);
    int big_mul_safe(char *a, char *b, char *ret);
    int big_exp_safe(char *a, char *b, char *ret);
#ifdef __cplusplus
} /// extern "C"
#endif
