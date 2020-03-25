#ifndef __PLAYER_H_
#define __PLAYER_H_

#include "./miniaudio/miniaudio.h"
#include <stdio.h>
#include <pthread.h>

#define DEBUG 0

#endif

static int verbose = DEBUG; // debug

typedef struct _mr_player mr_player;
typedef struct _mr_userdata mr_userdata;

typedef void(*callback) (void *pw);
//typedef void(*cb) (void* n);

struct _mr_userdata {
    ma_decoder *decoder;
    mr_player *player;
    int shielding_callback;

    uint64_t total_frame;
    callback cb;
    void *player_worker;
};

struct _mr_player {
    int exit;
    ma_device device;
    mr_userdata userdata;
};

char *mr_player_init(mr_player *p, ma_decoder* decoder, callback cb, void *pw); // 初始化
char *mr_player_start(mr_player *p); // 播放
void mr_player_stop(mr_player *p); // 停止设备播放
void mr_player_destory(mr_player *p); // 销毁
void mr_player_reset(mr_player *p); // 输入reset

void mr_curr_audio_info(mr_player *p, uint32_t *second, uint32_t *curr, uint32_t *sampleRate);

void data_callback(ma_device* pDevice, void* pOutput, const void* pInput, ma_uint32 frameCount);

// decoder
char *mr_decoder_init_file(ma_decoder* decoder, const char *filepath);
