'use strict';

const records = require( './mcc-mnc-list.json' );
const statusCodeList = require( './status-codes.json' );
const regionList = require('./regions.json');
const regionNames = regionList.map(function(x) { return x.name });

function all () {
  return records;
}

function statusCodes () {
  return statusCodeList;
}

function regions () {
  return regionList;
}

function filter ( filters ) {
  if (filters === undefined || filters === null) {
    return records;
  }

  if (typeof filters !== 'object') {
    throw new TypeError('Invalid parameter (object expected)');
  }

  let statusCode, mcc, mnc, countryCode, plmn, nibbledPlmn, region;

  if (filters.statusCode) {
    statusCode = filters.statusCode;
    if (statusCodeList.indexOf(statusCode) === -1) {
      throw new TypeError('Invalid statusCode parameter (not found in statusCode list)');
    }
  }

  if (filters.plmn) {
    if (typeof filters.plmn === 'string') {
      plmn = String(filters.plmn);
    } else {
      throw new TypeError('Invalid plmn parameter (string expected)');
    }
    mcc = plmn.substr(0, 3);
    mnc = plmn.substr(3);
  }

  if (filters.mcc && mcc) {
    throw new TypeError('Don\'t use mccmnc and mcc parameter at once');
  }
  if (filters.mnc && mnc) {
    throw new TypeError('Don\'t use mccmnc and mnc parameter at once');
  }

  if (filters.mcc) {
    if (typeof filters.mcc === 'string' || typeof filters.mcc === 'number') {
      mcc = String(filters.mcc);
    } else {
      throw new TypeError('Invalid mcc parameter (string expected)');
    }
  }

  if (filters.mnc) {
    if (typeof filters.mnc === 'string' || typeof filters.mnc === 'number') {
      mnc = String(filters.mnc);
    } else {
      throw new TypeError('Invalid mnc parameter (string expected)');
    }
  }
  
  if (filters.countryCode != undefined) {
    if (typeof filters.countryCode === 'string') {
      countryCode = filters.countryCode;
    } else {
      throw new TypeError('Invalid countryCode parameter (string expected)');
    }
  }

  if (filters.region) {
    region = filters.region;

    if (regionNames.indexOf(region) === -1) {
      throw new TypeError('Invalid region parameter (not found in region list)');
    }
  }

  let result = records;

  if (statusCode) {
    result = result.filter( record => record['status'] === statusCode );
  }
  if (countryCode) {
    result = result.filter( record => record['countryCode'] === countryCode );
  }
  if (region) {
    result = result.filter( record => record['region'] === region );
  }
  if (mcc) {
    result = result.filter( record => record['mcc'] === mcc );
  }
  if (mnc) {
    result = result.filter( record => record['mnc'] === mnc );
  }
  
  return result;
}

module.exports = { all, statusCodes, regions, filter };
