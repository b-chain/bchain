#include "sysApi.h"

class token
{
public:
    token() {}
    void create(char *sybmol, char *name, char *decimals, char *supply, int isIssue, unsigned long long blkNumber, unsigned int expiry);
    void transfer(char *to, char *amount, char *symbol, char *memo, unsigned long long blkNumber, unsigned int expiry);
    void issue(char *symbol, char *amount, char *memo, unsigned long long blkNumber, unsigned int expiry);
    void balanceOf(char *addr, char *symbol);
    void getSupply(char *symbol);
    void getDecimals(char *symbol);
    void getSymbol(char *symbol);
    void getName(char *symbol);
};

static void blkNumberValidate(unsigned long long blkNumber, unsigned int expiry)
{
    unsigned long long CurNumber = block_number();
    assert(((CurNumber >= blkNumber) && (CurNumber <= blkNumber+expiry)), "action expired");
}

void token::create(char *symbol, char *name, char *decimals, char *supply, int isIssue, unsigned long long blkNumber, unsigned int expiry)
{
    blkNumberValidate(blkNumber, expiry);
    char sender[48];
    action_sender(sender);
    requireAuth(sender);
    assert(strlen(symbol) < 64, "symbol len exceed");
    assert(strlen(name) < 64, "name len exceed");
    char key[128];
    int keyLen;

    keyLen = strjoint("symbol ", symbol, key);
    db_emplace(key, keyLen, symbol, strlen(symbol));

    keyLen = strjoint("isIssue ", symbol, key);
    db_emplace(key, keyLen, (char *)(&isIssue), 4);

    keyLen = strjoint("name ", symbol, key);
    db_emplace(key, keyLen, name, strlen(name));

    keyLen = strjoint("decimals ", symbol, key);
    db_emplace(key, keyLen, decimals, strlen(decimals));

    keyLen = strjoint("supply ", symbol, key);
    db_emplace(key, keyLen, supply, strlen(supply));

    char total_supply[128];
    int ts_len = big_exp_safe("10", decimals, total_supply);
    ts_len = big_mul_safe(total_supply, supply, total_supply);

    keyLen = strjoint(symbol, sender, key);
    db_emplace(key, keyLen, total_supply, ts_len);
}

void token::issue(char *symbol, char *amount, char *memo, unsigned long long blkNumber, unsigned int expiry)
{
    assert(strlen(memo) < 64, "issue memo exceed");
    blkNumberValidate(blkNumber, expiry);
    char sender[48];
    action_sender(sender);
    requireAuth(sender);
    assert(strlen(symbol) < 64, "symbol len exceed");

    char key[128];
    int keyLen;

    keyLen = strjoint("symbol ", symbol, key);
    int len = db_get(key, keyLen, symbol);
    assert(len > 0, "issue: symbol is not exist!");

    int isIssue = 0;
    keyLen = strjoint("isIssue ", symbol, key);
    db_get(key, keyLen, (char *)&isIssue);
    assert(isIssue != 0, "issue: token can not issue, illegal!");

    char supply[128];
    keyLen = strjoint("supply ", symbol, key);
    db_get(key, keyLen, supply);

    int sp_len = big_add(supply, amount, supply);
    db_set(key, keyLen, supply, sp_len);

    char decimals[128];
    keyLen = strjoint("decimals ", symbol, key);
    db_get(key, keyLen, decimals);

    char total_issue[128];
    big_exp_safe("10", decimals, total_issue);
    big_mul_safe(total_issue, amount, total_issue);

    char fromToken[128];
    keyLen = strjoint(symbol, sender, key);
    len = db_get(key, keyLen, fromToken);
    fromToken[len] = 0;
    int ft_len = big_add(fromToken, total_issue, fromToken);
    db_set(key, keyLen, fromToken, ft_len);
}
void token::transfer(char *to, char *amount, char *symbol, char *memo, unsigned long long blkNumber, unsigned int expiry)
{
    assert(strlen(memo) < 64, "transfer memo exceed");
    blkNumberValidate(blkNumber, expiry);
    char sender[48];
    action_sender(sender);
    char key[128];
    int keyLen;
    keyLen = strjoint("symbol ", symbol, key);
    int len = db_get(key, keyLen, symbol);
    assert(len > 0, "transfer: symbol is not exist!");

    char fromToken[128];
    keyLen = strjoint(symbol, sender, key);
    len = db_get(key, keyLen, fromToken);
    fromToken[len] = 0;
    int ft_len = big_sub_safe(fromToken, amount, fromToken);
    db_set(key, keyLen, fromToken, ft_len);

    char toToken[128];
    char key_to[128];
    int keyLen_to;
    str2lower(to);
    keyLen_to = strjoint(symbol, to, key_to);
    len = db_get(key_to, keyLen_to, toToken);
    toToken[len] = 0;
    int tt_len = big_add(toToken, amount, toToken);
    db_set(key_to, keyLen_to, toToken, tt_len);
}
void token::balanceOf(char *addr, char *symbol)
{
    char val[128];
    char key[128];
    int keyLen;
    str2lower(addr);
    keyLen = strjoint(symbol, addr, key);
    int len = db_get(key, keyLen, val);
    if (len == 0) {
        val[0] = '0';
        setResult(val, 1);
    } else {
        setResult(val, len);
    }
}

void token::getSupply(char *symbol)
{
    char val[128];
    char key[128];
    int keyLen;
    keyLen = strjoint("supply ", symbol, key);
    int len = db_get(key, keyLen, val);
    if (len == 0) {
        val[0] = '0';
        setResult(val, 1);
    } else {
        setResult(val, len);
    }
}

void token::getDecimals(char *symbol)
{
    char val[128];
    char key[128];
    int keyLen;
    keyLen = strjoint("decimals ", symbol, key);
    int len = db_get(key, keyLen, val);
    if (len == 0) {
        val[0] = 0;
        setResult(val, 1);
    } else {
        setResult(val, len);
    }
}

void token::getSymbol(char *symbol)
{
    char val[128];
    char key[128];
    int keyLen;
    keyLen = strjoint("symbol ", symbol, key);
    int len = db_get(key, keyLen, val);
    if (len == 0) {
        val[0] = 0;
        setResult(val, 1);
    } else {
        setResult(val, len);
    }
}

void token::getName(char *symbol)
{
    char val[128];
    char key[128];
    int keyLen;
    keyLen = strjoint("name ", symbol, key);
    int len = db_get(key, keyLen, val);
    if (len == 0) {
        val[0] = 0;
        setResult(val, 1);
    } else {
        setResult(val, len);
    }
}

extern "C"
{
    static token tk;
    void create(char *symbol, char *name, char *decimals, char *supply, int isIssue, unsigned long long blkNumber, unsigned int expiry)
    {
        return tk.create(symbol, name, decimals, supply, isIssue, blkNumber, expiry);
    }

    void issue(char *symbol, char *amount, char *memo, unsigned long long blkNumber, unsigned int expiry)
    {
        return tk.issue(symbol, amount, memo, blkNumber, expiry);
    }

    void transfer(char *to, char *amount, char *symbol, char *memo, unsigned long long blkNumber, unsigned int expiry)
    {
        return tk.transfer(to, amount, symbol, memo, blkNumber, expiry);
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
BCHAINIO_ABI(token, (create)(issue)(transfer)(balanceOf)(getSupply)(getDecimals)(getSymbol)(getName))