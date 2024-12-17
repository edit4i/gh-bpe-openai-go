use std::ffi::{c_char, CStr, CString};
use std::ptr;
use std::slice;

use bpe_openai::{cl100k_base, o200k_base, Tokenizer};

/// Opaque handle to a BPE tokenizer
#[repr(C)]
#[derive(Debug)]
pub struct TokenizerHandle {
    tokenizer: *const Tokenizer,
}

impl TokenizerHandle {
    fn new(tokenizer: &'static Tokenizer) -> *mut Self {
        let handle = Box::new(TokenizerHandle {
            tokenizer: tokenizer as *const Tokenizer,
        });
        Box::into_raw(handle)
    }

    unsafe fn get_tokenizer(handle: *const TokenizerHandle) -> Option<&'static Tokenizer> {
        if handle.is_null() {
            return None;
        }
        let handle = &*handle;
        Some(&*handle.tokenizer)
    }
}

#[no_mangle]
pub extern "C" fn bpe_cl100k_base() -> *mut TokenizerHandle {
    TokenizerHandle::new(cl100k_base())
}

#[no_mangle]
pub extern "C" fn bpe_o200k_base() -> *mut TokenizerHandle {
    TokenizerHandle::new(o200k_base())
}

#[no_mangle]
pub extern "C" fn bpe_count(handle: *const TokenizerHandle, text: *const c_char) -> usize {
    let tokenizer = match unsafe { TokenizerHandle::get_tokenizer(handle) } {
        Some(t) => t,
        None => return 0,
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
    handle: *const TokenizerHandle,
    text: *const c_char,
    limit: usize,
) -> usize {
    let tokenizer = match unsafe { TokenizerHandle::get_tokenizer(handle) } {
        Some(t) => t,
        None => return usize::MAX,
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
    handle: *const TokenizerHandle,
    text: *const c_char,
    token_count: *mut usize,
) -> *mut u32 {
    let tokenizer = match unsafe { TokenizerHandle::get_tokenizer(handle) } {
        Some(t) => t,
        None => return ptr::null_mut(),
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

    let mut tokens_vec = tokens.into_boxed_slice();
    let ptr = tokens_vec.as_mut_ptr();
    // Leak the box intentionally - it will be freed by Go
    let _ = Box::into_raw(tokens_vec);
    ptr
}

#[no_mangle]
pub extern "C" fn bpe_decode(
    handle: *const TokenizerHandle,
    tokens: *const u32,
    token_count: usize,
) -> *mut c_char {
    let tokenizer = match unsafe { TokenizerHandle::get_tokenizer(handle) } {
        Some(t) => t,
        None => return ptr::null_mut(),
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
pub extern "C" fn bpe_free(handle: *mut TokenizerHandle) {
    if !handle.is_null() {
        unsafe {
            let _handle = Box::from_raw(handle);
            // Don't drop the tokenizer since it's a static reference
        }
    }
}
