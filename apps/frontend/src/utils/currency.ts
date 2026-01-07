/**
 * Currency configuration utility
 * 
 * Provides centralized currency formatting using Intl.NumberFormat.
 * Currency can be configured via VITE_CURRENCY environment variable.
 * Default: IDR (Indonesian Rupiah)
 */

interface CurrencyConfig {
  currency: string
  locale: string
  symbol: string
}

const CURRENCY_CONFIG: Record<string, CurrencyConfig> = {
  IDR: { currency: 'IDR', locale: 'id-ID', symbol: 'Rp' },
  USD: { currency: 'USD', locale: 'en-US', symbol: '$' },
  EUR: { currency: 'EUR', locale: 'de-DE', symbol: '€' },
  GBP: { currency: 'GBP', locale: 'en-GB', symbol: '£' },
  JPY: { currency: 'JPY', locale: 'ja-JP', symbol: '¥' },
  SGD: { currency: 'SGD', locale: 'en-SG', symbol: 'S$' },
  MYR: { currency: 'MYR', locale: 'ms-MY', symbol: 'RM' },
}

const currencyKey = (import.meta.env.VITE_CURRENCY as string) || 'IDR'
const config = CURRENCY_CONFIG[currencyKey] || CURRENCY_CONFIG.IDR

/**
 * Format a number as currency
 * @param amount - The amount to format
 * @param options - Optional formatting options
 * @returns Formatted currency string
 */
export const formatCurrency = (
  amount: number | undefined | null,
  options?: { 
    decimals?: number
    showSymbol?: boolean 
  }
): string => {
  const value = amount ?? 0
  const decimals = options?.decimals ?? 0
  const showSymbol = options?.showSymbol ?? true

  if (!showSymbol) {
    return new Intl.NumberFormat(config.locale, {
      minimumFractionDigits: decimals,
      maximumFractionDigits: decimals,
    }).format(value)
  }

  return new Intl.NumberFormat(config.locale, {
    style: 'currency',
    currency: config.currency,
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(value)
}

/**
 * Format a number as currency for chart tooltips (without currency symbol)
 * @param amount - The amount to format
 * @returns Formatted number string with currency symbol prefix
 */
export const formatChartCurrency = (amount: number | undefined | null): string => {
  const value = amount ?? 0
  return `${config.symbol}${new Intl.NumberFormat(config.locale, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value)}`
}

/** Current currency symbol */
export const currencySymbol = config.symbol

/** Current currency code (e.g., 'IDR', 'USD') */
export const currencyCode = config.currency

/** Current locale for the currency */
export const currencyLocale = config.locale
