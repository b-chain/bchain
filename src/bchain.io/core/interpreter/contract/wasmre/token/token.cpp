#include "sysApi.h"

class token
{
public:
    token() {}
    void create(char *sybmol, char *name, int decimals,
                unsigned long long supply);
    void transfer(char *to, unsigned long long amount, char *symbol,
                  char *memo);
    void balanceOf(char *addr, char *symbol);
    void getSupply(char *symbol);
    void getDecimals(char *symbol);
    void getSymbol(char *symbol);
    void getName(char *symbol);
};

void token::create(char *symbol, char *name, int decimals,
                   unsigned long long supply)
{
    char sender[48];
    action_sender(sender);
    requireAuth(sender);
    assert(strlen(symbol) < 64, "symbol len exceed");
    assert(strlen(name) < 64, "name len exceed");
    char key[128];
    int keyLen;

    keyLen = strjoint("symbol ", symbol, key);
    db_emplace(key, keyLen, symbol, strlen(symbol));

    keyLen = strjoint("name ", symbol, key);
    db_emplace(key, keyLen, name, strlen(name));

    keyLen = strjoint("decimals ", symbol, key);
    db_emplace(key, keyLen, (char *)(&decimals), 4);

    keyLen = strjoint("supply ", symbol, key);
    db_emplace(key, keyLen, (char *)(&supply), 8);

    unsigned long long total_supply = supply;
    unsigned char *most = (unsigned char *)(&total_supply) + 7;
    assert(*most <= 24, "total supply exceed");
    for (int i = 0; i < decimals; i++)
    {
        total_supply *= 10;
        assert(*most <= 24, "total supply exceed!");
    }
    keyLen = strjoint(symbol, sender, key);
    db_emplace(key, keyLen, (char *)(&total_supply), 8);
}
void token::transfer(char *to, unsigned long long amount, char *symbol,
                     char *memo)
{
    char sender[48];
    action_sender(sender);
    char key[128];
    int keyLen;
    keyLen = strjoint("symbol ", symbol, key);
    int len = db_get(key, keyLen, symbol);
    assert(len > 0, "transfer: symbol is not exist!");

    unsigned long long fromToken = 0;
    keyLen = strjoint(symbol, sender, key);
    len = db_get(key, keyLen, (char *)&fromToken);
    assert(len == 8, "get sender token error");
    assert(fromToken >= amount, "insufficient token");
    fromToken -= amount;
    db_set(key, keyLen, (char *)(&fromToken), 8);

    unsigned long long toToken = 0;
    char key_to[128];
    int keyLen_to;
    str2lower(to);
    keyLen_to = strjoint(symbol, to, key_to);
    db_get(key_to, keyLen_to, (char *)&toToken);
    toToken += amount;
    db_set(key_to, keyLen_to, (char *)(&toToken), 8);
}
void token::balanceOf(char *addr, char *symbol)
{
    unsigned long long val = 0;
    char key[128];
    int keyLen;
    str2lower(addr);
    keyLen = strjoint(symbol, addr, key);
    db_get(key, keyLen, (char *)&val);
    setResult((char *)&val, 8);
}

void token::getSupply(char *symbol)
{
    unsigned long long supply = 0;
    char key[128];
    int keyLen;
    keyLen = strjoint("supply ", symbol, key);
    int len = db_get(key, keyLen, (char *)&supply);
    assert(len == 8, "get supply db error");
    setResult((char *)&supply, 8);
}

void token::getDecimals(char *symbol)
{
    int decimals = 0;
    char key[128];
    int keyLen;
    keyLen = strjoint("decimals ", symbol, key);
    int len = db_get(key, keyLen, (char *)&decimals);
    assert(len == 4, "get decimals db error");
    setResult((char *)&decimals, 4);
}

void token::getSymbol(char *symbol)
{
    char s[64];
    char key[128];
    int keyLen;
    keyLen = strjoint("symbol ", symbol, key);
    int len = db_get(key, keyLen, s);
    assert(len <= 64, "get symbol db error");
    setResult(s, len);
}

void token::getName(char *symbol)
{
    char name[64];
    char key[128];
    int keyLen;
    keyLen = strjoint("name ", symbol, key);
    int len = db_get(key, keyLen, name);
    assert(len <= 64, "get name db error");
    setResult(name, len);
}

extern "C"
{
    static token tk;
    void create(char *symbol, char *name, int decimals, unsigned long long supply)
    {
        return tk.create(symbol, name, decimals, supply);
    }

    void transfer(char *to, unsigned long long amount, char *symbol, char *memo)
    {
        return tk.transfer(to, amount, symbol, memo);
    }

    void balanceOf(char *addr, char *symbol)
    {
        return tk.balanceOf(addr, symbol);
    }

    void getSupply(char *symbol)
    {
        return tk.getSupply(symbol);
    }

    void getDecimals(char *symbol)
    {
        return tk.getDecimals(symbol);
    }

    void getSymbol(char *symbol)
    {
        return tk.getSymbol(symbol);
    }

    void getName(char *symbol)
    {
        return tk.getName(symbol);
    }
}

#define BCHAINIO_ABI(type, name)
// this macro is used for ABI generation declaration
BCHAINIO_ABI(token, (create)(balenceOf)(transfer)(getSupply)(getDecimals)(getSymbol)(getName))
