/**
 * Autogenerated by Thrift Compiler (0.13.0)
 *
 * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
 *  @generated
 */
package org.newrelic.nrjmx.v2.nrprotocol;


@javax.annotation.Generated(value = "Autogenerated by Thrift Compiler (0.13.0)", date = "2021-11-22")
public enum ValueType implements org.apache.thrift.TEnum {
  STRING(1),
  DOUBLE(2),
  INT(3),
  BOOL(4);

  private final int value;

  private ValueType(int value) {
    this.value = value;
  }

  /**
   * Get the integer value of this enum value, as defined in the Thrift IDL.
   */
  public int getValue() {
    return value;
  }

  /**
   * Find a the enum type by its integer value, as defined in the Thrift IDL.
   * @return null if the value is not found.
   */
  @org.apache.thrift.annotation.Nullable
  public static ValueType findByValue(int value) { 
    switch (value) {
      case 1:
        return STRING;
      case 2:
        return DOUBLE;
      case 3:
        return INT;
      case 4:
        return BOOL;
      default:
        return null;
    }
  }
}
