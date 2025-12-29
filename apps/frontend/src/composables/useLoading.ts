import { ref, computed } from 'vue'
import { useLoadingStore } from '@/stores/loading'

export function useLoading(defaultKey?: string) {
    const loadingStore = useLoadingStore()
    const localLoading = ref(false)

    // Use store-based loading if key provided, otherwise use local state
    const isLoading = computed(() => {
        if (defaultKey) {
            return loadingStore.isLoading(defaultKey)
        }
        return localLoading.value
    })

    const setLoading = (loading: boolean, key?: string) => {
        const targetKey = key || defaultKey
        if (targetKey) {
            loadingStore.setLoading(targetKey, loading)
        } else {
            localLoading.value = loading
        }
    }

    const withLoading = async <T>(
        operation: () => Promise<T>,
        key?: string
    ): Promise<T> => {
        const targetKey = key || defaultKey
        if (targetKey) {
            return loadingStore.withLoading(targetKey, operation)
        } else {
            try {
                localLoading.value = true
                return await operation()
            } finally {
                localLoading.value = false
            }
        }
    }

    // Convenience methods for common operations
    const executeWithLoading = async <T>(
        operation: () => Promise<T>,
        options: {
            key?: string
            onSuccess?: (result: T) => void
            onError?: (error: any) => void
            onFinally?: () => void
        } = {}
    ): Promise<T | null> => {
        try {
            const result = await withLoading(operation, options.key)

            if (options.onSuccess) {
                options.onSuccess(result)
            }

            return result
        } catch (error) {
            if (options.onError) {
                options.onError(error)
            }
            return null
        } finally {
            if (options.onFinally) {
                options.onFinally()
            }
        }
    }

    return {
        isLoading,
        setLoading,
        withLoading,
        executeWithLoading,
        // Expose loading keys for convenience
        LOADING_KEYS: loadingStore.LOADING_KEYS
    }
}

// Specialized hooks for common use cases
export function useAsyncOperation<T>(
    operation: () => Promise<T>,
    options: {
        immediate?: boolean
        loadingKey?: string
        onSuccess?: (result: T) => void
        onError?: (error: any) => void
    } = {}
) {
    const { isLoading, executeWithLoading } = useLoading(options.loadingKey)
    const data = ref<T | null>(null)
    const error = ref<any>(null)

    const execute = async () => {
        error.value = null

        const result = await executeWithLoading(operation, {
            key: options.loadingKey,
            onSuccess: (result) => {
                data.value = result
                if (options.onSuccess) {
                    options.onSuccess(result)
                }
            },
            onError: (err) => {
                error.value = err
                if (options.onError) {
                    options.onError(err)
                }
            }
        })

        return result
    }

    // Execute immediately if requested
    if (options.immediate) {
        execute()
    }

    return {
        data,
        error,
        isLoading,
        execute,
        refresh: execute
    }
}

// Hook for form submissions
export function useFormSubmission<T>(
    submitFn: () => Promise<T>,
    options: {
        loadingKey?: string
        onSuccess?: (result: T) => void
        onError?: (error: any) => void
        resetOnSuccess?: boolean
    } = {}
) {
    const { isLoading, executeWithLoading } = useLoading(options.loadingKey)
    const isSubmitted = ref(false)
    const submitError = ref<any>(null)

    const submit = async () => {
        submitError.value = null

        const result = await executeWithLoading(submitFn, {
            key: options.loadingKey,
            onSuccess: (result) => {
                isSubmitted.value = true
                if (options.onSuccess) {
                    options.onSuccess(result)
                }
                if (options.resetOnSuccess) {
                    setTimeout(() => {
                        isSubmitted.value = false
                    }, 2000)
                }
            },
            onError: (error) => {
                submitError.value = error
                if (options.onError) {
                    options.onError(error)
                }
            }
        })

        return result
    }

    const reset = () => {
        isSubmitted.value = false
        submitError.value = null
    }

    return {
        isLoading,
        isSubmitted,
        submitError,
        submit,
        reset
    }
}