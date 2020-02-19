import { environment } from '../../environments/environment';

export function determineLevel(currenStoage, type: 'adult' | 'child') {
  const baseNum =
    type === 'adult' ? environment.adultStorage : environment.childStorage;
  const ratio = (currenStoage / baseNum) * 100;
  return ratio >= 50
    ? 'safe'
    : ratio >= 20
    ? 'warning'
    : ratio > 0
    ? 'low'
    : 'soldout';
}

export function getlevelCSSColor(level) {
  return {
    safe: '48c774',
    warning: 'ffdd57',
    low: 'fc82b1',
    soldout: 'b9b9b9'
  }[level];
}

export function getlevelSize(level) {
  return {
    safe: 60,
    warning: 60,
    low: 50,
    soldout: 50
  }[level];
}

export const maskSortRule = maskOption => (a, b) => {
  return maskOption === 1
    ? a.maskAdult > b.maskAdult
      ? -1
      : 1
    : a.maskChild > b.maskChild
    ? -1
    : 1;
};

export const filterDataRule = filter => value => {
  return filter.length > 0
    ? [
        value.name.includes(filter),
        value.address.includes(filter),
        value.phone.includes(filter)
      ].some(r => r)
    : true;
};
