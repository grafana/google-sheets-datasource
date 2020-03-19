import { getGoogleSheetRangeInfoFromURL, formatCacheTimeLabel } from './QueryEditor';

describe('QueryEditor', () => {
  it('should extract id from URL', () => {
    const url = 'https://docs.google.com/spreadsheets/d/1m2idieRUHdzWTu3_cpYs1lUfP_jwfgL8NBaLtqLmia8/edit#gid=790763898&range=B19:F20';
    const info = getGoogleSheetRangeInfoFromURL(url);
    expect(info.spreadsheet).toBe('1m2idieRUHdzWTu3_cpYs1lUfP_jwfgL8NBaLtqLmia8');
    expect(info.range).toBe('B19:F20');
  });

  it('should format cache time seconds label correctly', () => {
    expect(formatCacheTimeLabel(0)).toBe('0s');
    expect(formatCacheTimeLabel(20)).toBe('20s');
    expect(formatCacheTimeLabel(60)).toBe('1m');
    expect(formatCacheTimeLabel(60 * 30)).toBe('30m');
    expect(formatCacheTimeLabel(60 * 60)).toBe('1h');
    expect(formatCacheTimeLabel(60 * 60 * 10)).toBe('10h');
  });
});
