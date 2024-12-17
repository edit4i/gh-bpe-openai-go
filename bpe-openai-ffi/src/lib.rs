use std::ffi::{c_char, CStr, CString};
use std::ptr;
use std::slice;

use bpe_openai::{cl100k_base, o200k_base, Tokenizer};

#[repr(C)]
pub struct bpe_tokenizer_t(Tokenizer);

#[no_mangle]
pub extern "C" fn bpe_cl100k_base() -> *mut bpe_tokenizer_t {
    Box::into_raw(Box::new(bpe_tokenizer_t(cl100k_base().clone())))
}

#[no_mangle]
pub extern "C" fn bpe_o200k_base() -> *mut bpe_tokenizer_t {
    Box::into_raw(Box::new(bpe_tokenizer_t(o200k_base().clone())))
}

#[no_mangle]
pub extern "C" fn bpe_count(tokenizer: *const bpe_tokenizer_t, text: *const c_char) -> usize {
    let tokenizer = unsafe {
        if tokenizer.is_null() {
            return 0;
        }
        &(*tokenizer).0
    };

    let text = unsafe {
        if text.is_null() {
            return 0;
        }
        match CStr::from_ptr(text).to_str() {
            Ok(s) => s,
            Err(_) => return 0,
        }
    };

    tokenizer.count(text)
}

#[no_mangle]
pub extern "C" fn bpe_count_till_limit(
    tokenizer: *const bpe_tokenizer_t,
    text: *const c_char,
    limit: usize,
) -> usize {
    let tokenizer = unsafe {
        if tokenizer.is_null() {
            return usize::MAX;
        }
        &(*tokenizer).0
    };

    let text = unsafe {
        if text.is_null() {
            return usize::MAX;
        }
        match CStr::from_ptr(text).to_str() {
            Ok(s) => s,
            Err(_) => return usize::MAX,
        }
    };

    match tokenizer.count_till_limit(text, limit) {
        Some(count) => count,
        None => usize::MAX,
    }
}

#[no_mangle]
pub extern "C" fn bpe_encode(
    tokenizer: *const bpe_tokenizer_t,
    text: *const c_char,
    token_count: *mut usize,
) -> *mut u32 {
    let tokenizer = unsafe {
        if tokenizer.is_null() {
            return ptr::null_mut();
        }
        &(*tokenizer).0
    };

    let text = unsafe {
        if text.is_null() {
            return ptr::null_mut();
        }
        match CStr::from_ptr(text).to_str() {
            Ok(s) => s,
            Err(_) => return ptr::null_mut(),
        }
    };

    let tokens = tokenizer.encode(text);
    unsafe {
        if !token_count.is_null() {
            *token_count = tokens.len();
        }
    }

    if tokens.is_empty() {
        return ptr::null_mut();
    }

    let mut tokens_vec = tokens.into_boxed_slice();
    let ptr = tokens_vec.as_mut_ptr();
    Box::into_raw(tokens_vec);
    ptr
}

#[no_mangle]
pub extern "C" fn bpe_decode(
    tokenizer: *const bpe_tokenizer_t,
    tokens: *const u32,
    token_count: usize,
) -> *mut c_char {
    let tokenizer = unsafe {
        if tokenizer.is_null() {
            return ptr::null_mut();
        }
        &(*tokenizer).0
    };

    let tokens = unsafe {
        if tokens.is_null() {
            return ptr::null_mut();
        }
        slice::from_raw_parts(tokens, token_count)
    };

    match tokenizer.decode(tokens) {
        Some(text) => match CString::new(text) {
            Ok(c_str) => c_str.into_raw(),
            Err(_) => ptr::null_mut(),
        },
        None => ptr::null_mut(),
    }
}

#[no_mangle]
pub extern "C" fn bpe_free(tokenizer: *mut bpe_tokenizer_t) {
    if !tokenizer.is_null() {
        unsafe {
            drop(Box::from_raw(tokenizer));
        }
    }
}
