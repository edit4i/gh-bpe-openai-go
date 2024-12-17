/* Auto-generated header for bpe-openai bindings */

#ifndef BPE_OPENAI_H
#define BPE_OPENAI_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * Opaque handle to a BPE tokenizer
 */
typedef struct bpe_TokenizerHandle {
  uint8_t _private[0];
} bpe_TokenizerHandle;

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

struct bpe_TokenizerHandle *bpe_cl100k_base(void);

struct bpe_TokenizerHandle *bpe_o200k_base(void);

uintptr_t bpe_count(const struct bpe_TokenizerHandle *handle, const char *text);

uintptr_t bpe_count_till_limit(const struct bpe_TokenizerHandle *handle,
                               const char *text,
                               uintptr_t limit);

uint32_t *bpe_encode(const struct bpe_TokenizerHandle *handle,
                     const char *text,
                     uintptr_t *token_count);

char *bpe_decode(const struct bpe_TokenizerHandle *handle,
                 const uint32_t *tokens,
                 uintptr_t token_count);

void bpe_free(struct bpe_TokenizerHandle *handle);

#ifdef __cplusplus
}  // extern "C"
#endif  // __cplusplus

#endif  /* BPE_OPENAI_H */
