#include "sysApi.h"

class account
{
public:
    account() {}
    void set(char* addr, char* val, int val_len);
    void get(char* addr);
};

static char *g_prefixAcc = "acc";
void account::set(char* addr, char* val, int val_len)
{
    char sender[48];
    action_sender(sender);
    requireAuth(sender);
    isHexAddress(addr);
    str2lower(addr);
    char key[128];
    int keyLen;
    keyLen = strjoint(g_prefixAcc, addr, key);
    db_set(key, keyLen, val, val_len);
}

void account::get(char *addr)
{
    char val[128];
    char key[128];
    int keyLen;
    str2lower(addr);
    keyLen = strjoint(g_prefixAcc, addr, key);
    int len = db_get(key, keyLen, val);
    setResult(val, len);
}

extern "C"
{
    static account acc;
    void set(char* addr, char* val, int val_len)
    {
        return acc.set(addr, val, val_len);
    }

    void get(char *addr)
    {
        return acc.get(addr);
    }
}

#define BCHAINIO_ABI(type, name)
// this macro is used for ABI generation declaration
BCHAINIO_ABI(account, (set)(get))