export function detectCountry(): string {
  if (typeof navigator === 'undefined' || !navigator.language) {
    return 'US';
  }

  const lang = navigator.language.toLowerCase();
  if (lang.includes('gb') || lang.includes('uk')) {
    return 'GB';
  }
  if (lang.includes('au')) {
    return 'AU';
  }
  if (lang.includes('in') || lang.startsWith('hi')) {
    return 'IN';
  }
  return 'US';
}
