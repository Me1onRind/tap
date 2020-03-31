#define DR_FLAC_IMPLEMENTATION

#include "./miniaudio/extras/dr_flac.h"  /* Enables FLAC decoding. */
#define DR_MP3_IMPLEMENTATION
#include "./miniaudio/extras/dr_mp3.h"   /* Enables MP3 decoding. */
#define DR_WAV_IMPLEMENTATION
#include "./miniaudio/extras/dr_wav.h"   /* Enables WAV decoding. */

#define MINIAUDIO_IMPLEMENTATION

#include "player.h"
#include <stdio.h>

#define MR_PANIC(str) if(verbose>0){printf("%s\n", str);};return str;
#define LOG(str) if(verbose>0) {printf("%s\n", str);}

char *mr_player_init(mr_player *p, ma_decoder* decoder, callback cb, void *pw, float volume) {
    p->userdata.player = p;
    p->userdata.decoder = decoder;

    p->userdata.cb = cb;
    p->userdata.player_worker = pw;
    p->userdata.shielding_callback= 0;

    p->userdata.total_frame = ma_decoder_get_length_in_pcm_frames(decoder);
    p->userdata.read_frame = 0;
    p->userdata.just_change= 0;

    pthread_mutex_init(&(p->userdata.lock), NULL);

    ma_device_config deviceConfig;
    deviceConfig = ma_device_config_init(ma_device_type_playback);
    deviceConfig.playback.format   = decoder->outputFormat;
    deviceConfig.playback.channels = decoder->outputChannels;
    deviceConfig.sampleRate        = decoder->outputSampleRate;
    deviceConfig.dataCallback      = data_callback;
    deviceConfig.pUserData         = &(p->userdata);

    if (ma_device_init(NULL, &deviceConfig, &(p->device)) != MA_SUCCESS) {
        MR_PANIC("Failed to open playback device.");
    }
    ma_device_set_master_volume(&(p->device), volume);

    LOG("player init success.");
    return NULL;
}

void mr_player_stop(mr_player *p) {
    if (ma_device_is_started(&(p->device))) {
        ma_device_stop(&(p->device));
    }
}

void mr_player_destory(mr_player *p) {
    ma_device_uninit(&(p->device));
    ma_decoder_uninit(p->userdata.decoder);
}

char *mr_player_start(mr_player *p) {
    if (ma_device_start(&(p->device)) != MA_SUCCESS) {
        MR_PANIC("Failed to start playback device.");
    }
    p->userdata.shielding_callback = 0;
    return NULL;
}


void mr_player_reset(mr_player *p) {
    if (ma_device_is_started(&(p->device))) {
        ma_device_stop(&(p->device));
    }
    pthread_mutex_lock(&(p->userdata.lock));
    p->userdata.read_frame = 0;
    ma_decoder_seek_to_pcm_frame(p->userdata.decoder, 0);
    pthread_mutex_unlock(&(p->userdata.lock));
}

void mr_curr_audio_info(mr_player *p, int64_t *second, int64_t *curr) {
    uint32_t sampleRate = p->userdata.decoder->outputSampleRate;
    *second = p->userdata.total_frame / sampleRate;
    pthread_mutex_lock(&(p->userdata.lock));
    *curr = *second * p->userdata.read_frame / p->userdata.total_frame;
    pthread_mutex_unlock(&(p->userdata.lock));
}

void mr_player_set_volume(mr_player *p, float volume) {
    ma_device_set_master_volume(&(p->device), volume);
}

void mr_player_seek_frame(mr_player *p, int64_t second) {
    pthread_mutex_lock(&(p->userdata.lock));
    /*int64_t curr = (int64_t)p->userdata.read_frame;*/
    int64_t sampleRate = (int64_t)p->userdata.decoder->outputSampleRate;
    int64_t position = sampleRate * (second + 1);
    if (position >=  p->userdata.total_frame) {
        position = p->userdata.total_frame;
    }
    if (position < 0) {
        position = 0;
    }
    if (ma_decoder_seek_to_pcm_frame(p->userdata.decoder, position) == MA_SUCCESS) {
        p->userdata.read_frame = position;
        p->userdata.just_change = 1;
    }
    pthread_mutex_unlock(&(p->userdata.lock));
}

void data_callback(ma_device* p_device, void* p_output, const void* p_input, ma_uint32 frame_count) {
    mr_userdata* userdata = (mr_userdata*)p_device->pUserData;
    if (userdata == NULL) {
        return;
    }
    pthread_mutex_lock(&(userdata->lock));
    if (userdata->just_change == 1) {
        userdata->just_change = 0;
        pthread_mutex_unlock(&(userdata->lock));
        usleep(600 * 1000);
        return;
    }

    uint64_t n_read = ma_decoder_read_pcm_frames(userdata->decoder, p_output, frame_count);
    userdata->read_frame += n_read;

    if (n_read < frame_count && !userdata->shielding_callback) {
        if (userdata->cb) {
            userdata->cb(userdata->player_worker);
        }
        userdata->shielding_callback= 1;

    }
    pthread_mutex_unlock(&(userdata->lock));
    (void)p_input;
}


char *mr_decoder_init_file(ma_decoder *decoder, const char *filepath) {
    if (ma_decoder_init_file(filepath, NULL, decoder) != MA_SUCCESS) {
        return "decoder init failed";
    }
    return NULL;
}
